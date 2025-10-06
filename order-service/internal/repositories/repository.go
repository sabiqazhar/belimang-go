package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	intDB "github.com/sabiqazhar/belimang-go/order-service/internal/db"
)

type OrderRepository interface {
	InsertOrderItems(ctx context.Context, param []intDB.InsertOrderItemsParams) (int64, error)
	InsertOrder(ctx context.Context, param intDB.InsertOrderParams) (int64, error)
	WithTx(tx pgx.Tx) OrderRepository
}
