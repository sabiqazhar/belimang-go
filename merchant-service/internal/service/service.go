package service

import (
	"context"

	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
)

type MerchantService interface {
	CreateMerchant(ctx context.Context, req model.CreateMerchantRequest) (int64, error)
	GetMerchants(ctx context.Context, param db.GetMerchantListParams) ([]db.GetMerchantListRow, error)
	AddItem(ctx context.Context, req model.AddItemRequest, merchantId int64) (int64, error)
	GetItems(ctx context.Context, param db.GetItemListParams) ([]db.GetItemListRow, error)
}
