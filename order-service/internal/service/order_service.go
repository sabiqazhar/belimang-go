package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sabiqazhar/belimang-go/order-service/helper"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	"github.com/sabiqazhar/belimang-go/order-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type OrderServiceImpl struct {
	orderRepo repositories.OrderRepository
	store     *pgxpool.Pool
}

func NewOrderServiceImpl(orderRepo repositories.OrderRepository, dataStore *pgxpool.Pool) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		store:     dataStore,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req model.CreateEstimateRequest, userId int32) (int64, error) {
	var orderId int64

	if err := s.validateStartingPoint(req); err != nil {
		return 0, err
	}

	tx, err := s.store.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	txRepo := s.orderRepo.WithTx(tx)
	amount := helper.FloatToNumeric(100.00)
	if amount == (pgtype.Numeric{}) {
		return 0, errors.New("amount is zero")
	}

	order := db.InsertOrderParams{
		CustomerID:                     userId,
		Status:                         "pending",
		OrderDate:                      pgtype.Timestamp{Time: time.Now()},
		Longitude:                      pgtype.Float8{Float64: req.UserLocation.Long, Valid: true},
		Latitude:                       pgtype.Float8{Float64: req.UserLocation.Lat, Valid: true},
		TotalAmount:                    amount, // Example fixed amount
		EstimatedDeliveryTimeInMinutes: 400,    // Example fixed time
	}

	orderID, err := txRepo.InsertOrder(ctx, order)
	orderId = orderID
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return 0, err
	}

	var orderItems []db.InsertOrderItemsParams
	for _, orderRow := range req.Orders {
		for _, item := range orderRow.Items {
			orderItems = append(orderItems, db.InsertOrderItemsParams{
				MerchantID:    orderRow.MerchantID,
				OrderID:       pgtype.Int8{Int64: orderID, Valid: true},
				ProductID:     item.ItemID,
				Quantity:      int32(item.Quantity),
				Price:         helper.FloatToNumeric(200),
				StartingPoint: pgtype.Bool{Bool: *orderRow.IsStartingPoint, Valid: true},
			})
		}
	}

	_, err = txRepo.InsertOrderItems(ctx, orderItems)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return orderId, nil
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
