package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/order-service/helper"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	"github.com/sabiqazhar/belimang-go/order-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type OrderServiceImpl struct {
	orderRepo repositories.OrderRepository
	store     *repositories.Store
}

func NewOrderServiceImpl(orderRepo repositories.OrderRepository, dataStore *repositories.Store) OrderService {
	return &OrderServiceImpl{
		orderRepo: orderRepo,
		store:     dataStore,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req model.CreateEstimateRequest, userId int32) (int64, error) {
	var orderId int64
	err := s.store.ExecTx(ctx, func(queries *db.Queries) error {
		amount := helper.FloatToNumeric(100.00)
		if amount == (pgtype.Numeric{}) {
			return errors.New("amount is zero")
		}

		order := db.InsertOrderParams{
			CustomerID:                     userId,
			Status:                         "pending",
			TotalAmount:                    amount, // Example fixed amount
			EstimatedDeliveryTimeInMinutes: 400,    // Example fixed time
		}

		orderID, err := queries.InsertOrder(ctx, order)
		orderId = orderID
		if err != nil {
			logger.Logger.Error().Err(err).Msg("failed to insert order")
			return err
		}

		var orderItems []db.InsertOrderItemsParams
		for _, orderRow := range req.Orders {
			for _, item := range orderRow.Items {
				orderItems = append(orderItems, db.InsertOrderItemsParams{
					MerchantID:    orderRow.MerchantID,
					OrderID:       pgtype.Int8{Int64: orderID, Valid: true},
					ProductID:     int32(item.ItemID),
					Quantity:      int32(item.Quantity),
					Price:         helper.FloatToNumeric(200),
					StartingPoint: pgtype.Bool{Bool: orderRow.IsStartingPoint, Valid: true},
				})
			}
		}

		_, err = queries.InsertOrderItems(ctx, orderItems)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("failed to insert order")
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return 0, err
	}

	return orderId, nil
}
