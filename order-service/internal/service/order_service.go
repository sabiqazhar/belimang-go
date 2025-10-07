package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/sabiqazhar/belimang-go/order-service/internal/client"
	pb "github.com/sabiqazhar/belimang-go/proto/merchant"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sabiqazhar/belimang-go/order-service/helper"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	"github.com/sabiqazhar/belimang-go/order-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type OrderServiceImpl struct {
	orderRepo      repositories.OrderRepository
	store          *pgxpool.Pool
	merchantClient *client.MerchantClient
}

func NewOrderServiceImpl(orderRepo repositories.OrderRepository, dataStore *pgxpool.Pool, merchantClient *client.MerchantClient) OrderService {
	return &OrderServiceImpl{
		orderRepo:      orderRepo,
		store:          dataStore,
		merchantClient: merchantClient,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req model.CreateEstimateRequest, userId int32) (model.CreateOrderResponse, error) {
	if err := s.validateStartingPoint(req); err != nil {
		return model.CreateOrderResponse{}, err
	}

	tx, err := s.store.Begin(ctx)
	if err != nil {
		return model.CreateOrderResponse{}, err
	}
	defer tx.Rollback(ctx)

	txRepo := s.orderRepo.WithTx(tx)

	orders, err := s.gatherMerchantDetailFromOrder(ctx, req.Orders)
	if err != nil {
		return model.CreateOrderResponse{}, err
	}

	estimatedDistance := s.calculateEstimate(orders, req.UserLocation)
	estimatedDeliveryTimeInMinutes := helper.CalculateDeliveryTime(estimatedDistance)
	defaultAmount := helper.FloatToNumeric(0)
	if defaultAmount == (pgtype.Numeric{}) {
		return model.CreateOrderResponse{}, errors.New("amount is zero")
	}

	order := db.InsertOrderParams{
		CustomerID:                     userId,
		Status:                         "pending",
		OrderDate:                      pgtype.Timestamp{Time: time.Now()},
		Longitude:                      pgtype.Float8{Float64: req.UserLocation.Long, Valid: true},
		Latitude:                       pgtype.Float8{Float64: req.UserLocation.Lat, Valid: true},
		TotalAmount:                    defaultAmount,
		EstimatedDeliveryTimeInMinutes: int32(estimatedDeliveryTimeInMinutes),
		TotalDistanceInMeters:          int32(estimatedDistance),
	}

	orderID, err := txRepo.InsertOrder(ctx, order)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return model.CreateOrderResponse{}, err
	}

	totalAmount, orderItems, err := s.constructInsertOrderItems(ctx, req, orderID)
	err = txRepo.UpdateOrderTotalAmount(ctx, db.UpdateOrderAmountParams{
		ID:          orderID,
		TotalAmount: helper.FloatToNumeric(totalAmount),
	})
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to update order amount")
		return model.CreateOrderResponse{}, err
	}

	_, err = txRepo.InsertOrderItems(ctx, orderItems)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return model.CreateOrderResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.CreateOrderResponse{}, err
	}

	return model.CreateOrderResponse{
		CalculatedEstimateId:           orderID,
		EstimatedDeliveryTimeInMinutes: estimatedDeliveryTimeInMinutes,
		TotalPrice:                     totalAmount,
	}, nil
}

func (s *OrderServiceImpl) validateStartingPoint(req model.CreateEstimateRequest) error {
	var isStartingPointSeen bool
	for _, orderRow := range req.Orders {
		if *orderRow.IsStartingPoint {
			if isStartingPointSeen {
				return errors.New("only one starting point is allowed")
			}
			isStartingPointSeen = true
		}
	}
	if !isStartingPointSeen {
		return errors.New("starting point is required")
	}
	return nil
}

func (s *OrderServiceImpl) getMerchantDetail(ctx context.Context, merchantID int64) (model.MerchantDetail, error) {
	merchantReq := &pb.GetMerchantRequest{
		Id: merchantID,
	}
	resp, err := s.merchantClient.Client.GetMerchant(ctx, merchantReq)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to get merchant detail")
		return model.MerchantDetail{}, err
	}

	if resp.Error {
		return model.MerchantDetail{}, errors.New("merchant not found")
	}
	return model.MerchantDetail{
		ID: strconv.FormatInt(resp.Id, 10),
		Location: model.LongLat{
			Lat:  resp.Latitude,
			Long: resp.Longitude,
		},
	}, nil
}

func (s *OrderServiceImpl) gatherMerchantDetailFromOrder(ctx context.Context, merchants []model.Order) ([]model.MerchantDetail, error) {
	var merchantDetails []model.MerchantDetail
	var startingPointMerchantIdx int64

	for idx, merchant := range merchants {
		merchantInt, err := strconv.ParseInt(merchant.MerchantID, 10, 64)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("failed to parse merchant id")
			return []model.MerchantDetail{}, err
		}

		merchantDetail, err := s.getMerchantDetail(ctx, merchantInt)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("failed to get merchant detail")
			return []model.MerchantDetail{}, err
		}

		merchantDetails = append(merchantDetails, merchantDetail)
		if *merchant.IsStartingPoint {
			startingPointMerchantIdx = int64(idx)
		}
	}
	merchantDetails[0], merchantDetails[startingPointMerchantIdx] = merchantDetails[startingPointMerchantIdx], merchantDetails[0]

	return merchantDetails, nil
}

func (s *OrderServiceImpl) calculateEstimate(merchants []model.MerchantDetail, finalLocation model.LongLat) float64 {
	var totalDistance float64
	for _, merchant := range merchants[:len(merchants)-1] {
		totalDistance += helper.CalculateHaversineDistance(merchant.Location.Lat, merchant.Location.Long, merchants[1].Location.Lat, merchants[1].Location.Long)
	}

	totalDistance += helper.CalculateHaversineDistance(merchants[len(merchants)-1].Location.Lat, merchants[len(merchants)-1].Location.Long, finalLocation.Lat, finalLocation.Long)
	return totalDistance
}

func (s *OrderServiceImpl) constructInsertOrderItems(ctx context.Context, req model.CreateEstimateRequest, orderID int64) (float64, []db.InsertOrderItemsParams, error) {
	var orderItems []db.InsertOrderItemsParams
	var totalAmount float64
	for _, orderRow := range req.Orders {
		intMerchantID, err := strconv.ParseInt(orderRow.MerchantID, 10, 64)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("failed to parse merchant id")
			return 0, []db.InsertOrderItemsParams{}, err
		}

		for _, item := range orderRow.Items {

			intItemID, err := strconv.ParseInt(item.ItemID, 10, 64)
			if err != nil {
				logger.Logger.Error().Err(err).Msg("failed to parse item id")
				return 0, []db.InsertOrderItemsParams{}, err
			}

			itemDetail, err := s.merchantClient.Client.GetItemDetail(ctx, &pb.GetItemDetailRequest{
				Id: intItemID,
			})

			if err != nil || itemDetail.Error {
				logger.Logger.Error().Err(err).Msg("failed to get item detail")
				return 0, []db.InsertOrderItemsParams{}, err
			}

			if itemDetail.MerchantId != intMerchantID {
				logger.Logger.Error().Err(err).Msg("merchant id does not match")
			}

			totalAmount += float64(item.Quantity) * itemDetail.Price
			orderItems = append(orderItems, db.InsertOrderItemsParams{
				MerchantID:    orderRow.MerchantID,
				OrderID:       pgtype.Int8{Int64: orderID, Valid: true},
				ProductID:     item.ItemID,
				Quantity:      int32(item.Quantity),
				Price:         helper.FloatToNumeric(itemDetail.Price),
				StartingPoint: pgtype.Bool{Bool: *orderRow.IsStartingPoint, Valid: true},
			})
		}
	}
	return totalAmount, orderItems, nil
}
