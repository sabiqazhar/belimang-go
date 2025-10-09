package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/order-service/helper"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	pb "github.com/sabiqazhar/belimang-go/proto/merchant"
)

type FetchMerchantJob struct {
	MerchantID      string
	Index           int // Important! To maintain order
	IsStartingPoint bool
}

type FetchMerchantResult struct {
	Detail          model.MerchantDetail
	Index           int // Match with job
	IsStartingPoint bool
	Error           error
}

func (s *OrderServiceImpl) gatherMerchantDetailConcurrent(ctx context.Context, merchants []model.Order) ([]model.MerchantDetail, error) {
	jobCount := len(merchants)

	jobs := make(chan FetchMerchantJob, jobCount)
	results := make(chan FetchMerchantResult, jobCount)

	workerCount := 5
	if jobCount < workerCount {
		workerCount = jobCount
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i <= workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.merchantFetchWorker(ctx, jobs, results)
		}()
	}

	//send jobs to channel
	for idx, merchant := range merchants {
		jobs <- FetchMerchantJob{
			MerchantID:      merchant.MerchantID,
			Index:           idx,
			IsStartingPoint: *merchant.IsStartingPoint,
		}
	}
	close(jobs) // Tell workers: no more jobs coming

	// Close results when all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from channel
	resultSlice := make([]FetchMerchantResult, 0, jobCount)
	for result := range results { // Blocks until channel closes
		resultSlice = append(resultSlice, result)
	}

	merchantDetails := make([]model.MerchantDetail, len(merchants))
	var startingPointIdx int
	var firstError error

	for _, result := range resultSlice {
		// Check for errors
		if result.Error != nil {
			if firstError == nil {
				firstError = result.Error
			}
			continue
		}

		// Put result in correct position
		merchantDetails[result.Index] = result.Detail

		// Track starting point
		if result.IsStartingPoint {
			startingPointIdx = result.Index
		}
	}

	// If any error occurred, return it
	if firstError != nil {
		return nil, firstError
	}

	// Swap starting point to position 0
	if startingPointIdx != 0 {
		merchantDetails[0], merchantDetails[startingPointIdx] =
			merchantDetails[startingPointIdx], merchantDetails[0]
	}

	return merchantDetails, nil

}

func (s *OrderServiceImpl) merchantFetchWorker(
	ctx context.Context,
	jobs <-chan FetchMerchantJob,
	results chan<- FetchMerchantResult,
) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Channel closed, no more jobs
				return
			}

			// Process this one job
			result := s.processMerchantJob(ctx, job)

			// Send result back
			results <- result

		case <-ctx.Done():
			// User cancelled request
			return
		}
	}
}

func (s *OrderServiceImpl) processMerchantJob(ctx context.Context, job FetchMerchantJob) FetchMerchantResult {
	// Step 1: Parse merchant ID
	merchantInt, err := strconv.ParseInt(job.MerchantID, 10, 64)
	if err != nil {
		// What should we return on error?
		return FetchMerchantResult{
			Index: job.Index, // Important! Keep index
			Error: fmt.Errorf("failed to parse merchant ID: %w", err),
		}
	}

	// Step 2: Fetch merchant detail (network call)
	merchantDetail, err := s.getMerchantDetail(ctx, merchantInt)
	if err != nil {
		return FetchMerchantResult{
			Index: job.Index,
			Error: fmt.Errorf("failed to get merchant detail: %w", err),
		}
	}

	// Step 3: Return success result
	return FetchMerchantResult{
		Detail:          merchantDetail,
		Index:           job.Index,
		IsStartingPoint: job.IsStartingPoint,
		Error:           nil,
	}
}

// concurrent for fetching item details
// What information does a worker need to process ONE item?
type OrderItem struct {
	ItemID   string
	Quantity int32
}

type FetchItemJob struct {
	OrderRow model.Order // Need merchant ID, isStartingPoint
	Item     OrderItem   // Need item ID, quantity
	OrderID  int64
	Index    int // To maintain order (important!)
}

// What result do we expect back?
type FetchItemResult struct {
	OrderItem db.InsertOrderItemsParams
	ItemPrice float64 // quantity × price
	Index     int
	Error     error
}

// ConstructInsertOrderItemsConcurrent
func (s *OrderServiceImpl) ConstructInsertOrderItemsConcurrent(
	ctx context.Context,
	req model.CreateEstimateRequest,
	orderID int64,
) (float64, []db.InsertOrderItemsParams, error) {

	// count number of items
	var totalItems int64
	for _, order := range req.Orders {
		totalItems += int64(len(order.Items))
	}

	if totalItems == 0 {
		return 0, []db.InsertOrderItemsParams{}, nil
	}

	jobs := make(chan FetchItemJob, totalItems)
	results := make(chan FetchItemResult, totalItems)

	workerNums := 10
	if totalItems < int64(workerNums) {
		workerNums = int(totalItems)
	}

	var wg sync.WaitGroup
	for i := 0; i < workerNums; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.itemFetchWorker(ctx, jobs, results)
		}()
	}

	jobIndex := 0
	for _, orderRow := range req.Orders {
		for _, item := range orderRow.Items {
			jobs <- FetchItemJob{
				OrderRow: orderRow,
				Item:     OrderItem{ItemID: item.ItemID, Quantity: int32(item.Quantity)},
				OrderID:  orderID,
				Index:    jobIndex,
			}
			jobIndex++
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	resultSlice := make([]FetchItemResult, 0, totalItems)
	for result := range results { // Blocks until channel closes
		resultSlice = append(resultSlice, result)
	}

	// Now results are in random order, need to sort by Index
	// Pre-allocate final slices
	orderItems := make([]db.InsertOrderItemsParams, totalItems)
	var totalAmount float64
	var firstError error

	// Process results in order
	for _, result := range resultSlice {
		// Check for errors
		if result.Error != nil {
			if firstError == nil {
				firstError = result.Error
			}
			continue // Skip this item but check others
		}

		// Put result in correct position
		orderItems[result.Index] = result.OrderItem
		totalAmount += result.ItemPrice
	}

	// If any error, return it
	if firstError != nil {
		return 0, nil, firstError
	}

	return totalAmount, orderItems, nil
}

func (s *OrderServiceImpl) itemFetchWorker(
	ctx context.Context,
	jobs <-chan FetchItemJob,
	results chan<- FetchItemResult,
) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				// Channel closed, no more jobs
				return
			}

			// Process this one job
			result := s.processItemJob(ctx, job)
			// Send result back
			results <- result

		case <-ctx.Done():
			// User cancelled request
			return
		}
	}
}

func (s *OrderServiceImpl) processItemJob(
	ctx context.Context,
	job FetchItemJob,
) FetchItemResult {

	// Step 1: Parse merchant ID
	intMerchantID, err := strconv.ParseInt(job.OrderRow.MerchantID, 10, 64)
	if err != nil {
		return FetchItemResult{
			Index: job.Index,
			Error: fmt.Errorf("failed to parse merchant id: %w", err),
		}
	}

	// Step 2: Parse item ID
	intItemID, err := strconv.ParseInt(job.Item.ItemID, 10, 64)
	if err != nil {
		return FetchItemResult{
			Index: job.Index,
			Error: fmt.Errorf("failed to parse item id: %w", err),
		}
	}

	// Step 3: Fetch item details (NETWORK CALL)
	itemDetail, err := s.merchantClient.Client.GetItemDetail(ctx, &pb.GetItemDetailRequest{
		Id: intItemID,
	})

	if err != nil || itemDetail.Error {
		return FetchItemResult{
			Index: job.Index,
			Error: fmt.Errorf("failed to get item detail: %w", err),
		}
	}

	// Step 4: Validate merchant matches
	if itemDetail.MerchantId != intMerchantID {
		return FetchItemResult{
			Index: job.Index,
			Error: errors.New("merchant id does not match"),
		}
	}

	// Step 5: Calculate price
	itemPrice := float64(job.Item.Quantity) * itemDetail.Price

	// Step 6: Build order item
	orderItem := db.InsertOrderItemsParams{
		MerchantID:    job.OrderRow.MerchantID,
		OrderID:       pgtype.Int8{Int64: job.OrderID, Valid: true},
		ProductID:     job.Item.ItemID,
		Quantity:      int32(job.Item.Quantity),
		Price:         helper.FloatToNumeric(itemDetail.Price),
		StartingPoint: pgtype.Bool{Bool: *job.OrderRow.IsStartingPoint, Valid: true},
	}

	// Step 7: Return success
	return FetchItemResult{
		OrderItem: orderItem,
		ItemPrice: itemPrice,
		Index:     job.Index,
		Error:     nil,
	}
}
