package service

import (
	"context"

	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
)

type MerchantService interface {
	CreateMerchant(ctx context.Context, req model.CreateMerchantRequest) (int64, error)
}
