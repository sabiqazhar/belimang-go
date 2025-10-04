package repositories

import (
	"context"
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/db"
)

type UploadPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewUploadRepository(database *pgx.Conn, nodeID int64) (UploadRepository, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	
	return &UploadPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}, nil
}

func (r *UploadPostgresRepo) UploadImage(ctx context.Context, image db.InsertImageParams) (int64, error) {
	id := r.snowflakeNode.Generate().Int64()
	image.ID = id
	_, err := r.db.InsertImage(ctx, image)

	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to create image upload data")
		return 0, err
	}

	return id, nil
}