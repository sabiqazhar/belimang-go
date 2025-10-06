package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	intDB "github.com/sabiqazhar/belimang-go/order-service/internal/db"
)

type OrderRepository interface {
	InsertOrderItems(ctx context.Context, param []intDB.InsertOrderItemsParams) (int64, error)
	InsertOrder(ctx context.Context, param intDB.InsertOrderParams) (int64, error)
}

type Store struct {
	*intDB.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: intDB.New(db),
	}
}

func (store *Store) ExecTx(ctx context.Context, fn func(*intDB.Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := intDB.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
