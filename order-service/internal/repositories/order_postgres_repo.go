package repositories

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type OrderPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewOrderRepository(database *pgx.Conn, node *snowflake.Node) OrderRepository {
	return &OrderPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}
}

func (r *OrderPostgresRepo) InsertOrder(ctx context.Context, param db.InsertOrderParams) (int64, error) {
	snowID := r.snowflakeNode.Generate()
	param.ID = snowID.Int64()

	orderID, err := r.db.InsertOrder(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order")
		return 0, err
	}

	return orderID, nil
}

func (r *OrderPostgresRepo) InsertOrderItems(ctx context.Context, param []db.InsertOrderItemsParams) (int64, error) {
	for _, item := range param {
		snowID := r.snowflakeNode.Generate().Int64()
		item.ID = snowID
	}

	rowsInserted, err := r.db.InsertOrderItems(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order items")
		return 0, err
	}
	return rowsInserted, nil
}
