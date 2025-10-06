package repositories

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sabiqazhar/belimang-go/order-service/internal/db"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type OrderPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewOrderRepository(database *pgxpool.Pool, node *snowflake.Node) OrderRepository {
	return &OrderPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}
}

func (r *OrderPostgresRepo) WithTx(tx pgx.Tx) OrderRepository {
	return &OrderPostgresRepo{
		db:            db.New(tx),
		snowflakeNode: r.snowflakeNode,
	}
}

func (r *OrderPostgresRepo) InsertOrder(ctx context.Context, param db.InsertOrderParams) (int64, error) {
	snowID := r.snowflakeNode.Generate()
	param.ID = snowID.Int64()
	logger.Logger.Info().Interface("param", param).Msg("param")
	orderID, err := r.db.InsertOrder(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("[error occurred on InsertOrder(ctx context.Context, param db.InsertOrderParams)]failed to insert order")
		return 0, err
	}

	return orderID, nil
}

func (r *OrderPostgresRepo) InsertOrderItems(ctx context.Context, param []db.InsertOrderItemsParams) (int64, error) {
	for i := range param {
		snowID := r.snowflakeNode.Generate().Int64()
		param[i].ID = snowID
	}

	logger.Logger.Info().Interface("orderItems", param).Msg("Inserting order items")
	rowsInserted, err := r.db.InsertOrderItems(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to insert order items")
		return 0, err
	}
	return rowsInserted, nil
}
