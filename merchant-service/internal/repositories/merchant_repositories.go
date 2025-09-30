package repositories

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
)

type MerchantPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewMerchantRepository(database *pgx.Conn, node *snowflake.Node) MerchantRepository {
	return &MerchantPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}
}
