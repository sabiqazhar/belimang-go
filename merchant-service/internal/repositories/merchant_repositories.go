package repositories

import (
	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	"golang.org/x/net/context"
)

type MerchantPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewMerchantRepository(database *pgxpool.Pool, node *snowflake.Node) MerchantRepository {
	return &MerchantPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}
}

func (m *MerchantPostgresRepo) InsertMerchant(ctx context.Context, merchant db.CreateMerchantParams) (db.Merchants, error) {
	merchant.ID = m.snowflakeNode.Generate().Int64()
	return m.db.CreateMerchant(ctx, merchant)
}

func (m *MerchantPostgresRepo) GetMerchants(ctx context.Context, merchantParam db.GetMerchantListParams) ([]db.GetMerchantListRow, error) {
	merchants, err := m.db.GetMerchantList(ctx, merchantParam)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("GetMerchants")
		return nil, err
	}
	return merchants, nil
}

func (m *MerchantPostgresRepo) AddItem(ctx context.Context, item db.AddItemParams) (int64, error) {
	item.ID = m.snowflakeNode.Generate().Int64()
	itemID, err := m.db.AddItem(ctx, item)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("AddItem")
		return 0, err
	}
	return itemID, nil
}

func (m *MerchantPostgresRepo) GetItems(ctx context.Context, param db.GetItemListParams) ([]db.GetItemListRow, error) {
	items, err := m.db.GetItemList(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("GetItem")
		return nil, err
	}
	return items, nil
}

func (m *MerchantPostgresRepo) GetMerchantById(ctx context.Context, id int64) (db.GetMerchantByIdRow, error) {
	merchant, err := m.db.GetMerchantById(ctx, id)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("GetMerchantById")
		return db.GetMerchantByIdRow{}, err
	}
	return merchant, nil
}

func (m *MerchantPostgresRepo) GetItemByID(ctx context.Context, id int64) (db.GetItemByIDRow, error) {
	itemDetail, err := m.db.GetItemByID(ctx, id)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("GetMerchantItemByID")
		return db.GetItemByIDRow{}, err
	}
	return itemDetail, nil
}
