package repositories

import (
	"context"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
)

type MerchantRepository interface {
	InsertMerchant(ctx context.Context, merchant db.CreateMerchantParams) (db.Merchants, error)
}
