package service

import (
	"context"

	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req model.CreateEstimateRequest, userId int32) (int64, error)
}
