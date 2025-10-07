package repositories

import (
	"context"

	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
)

type MerchantRepository interface {
	InsertMerchant(ctx context.Context, merchant db.CreateMerchantParams) (db.Merchants, error)
	GetMerchants(ctx context.Context, merchantParam db.GetMerchantListParams) ([]db.GetMerchantListRow, error)
	AddItem(ctx context.Context, item db.AddItemParams) (int64, error)
	GetItems(ctx context.Context, param db.GetItemListParams) ([]db.GetItemListRow, error)
	GetMerchantById(ctx context.Context, id int64) (db.GetMerchantByIdRow, error)
	GetItemByID(ctx context.Context, id int64) (db.GetItemByIDRow, error)
}
