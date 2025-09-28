package repositories

import (
	"context"
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	"github.com/sabiqazhar/belimang-go/user-service/internal/db"
)

type UserPostgresRepo struct {
	db            *db.Queries
	snowflakeNode *snowflake.Node
}

func NewUserRepository(database *pgx.Conn, nodeID int64) (UserRepository, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}

	return &UserPostgresRepo{
		db:            db.New(database),
		snowflakeNode: node,
	}, nil
}

func (s *UserPostgresRepo) CreateUser(ctx context.Context, user db.CreateUserParams) (int64, error) {
	id := s.snowflakeNode.Generate().Int64()
	user.ID = id
	_, err := s.db.CreateUser(ctx, user)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to create user")
		return 0, err
	}
	return id, nil
}

func (s *UserPostgresRepo) IsEmailAdminExists(ctx context.Context, email string) (bool, error) {
	user, err := s.db.IsAdminEmailExists(ctx, email)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to validate user admin")
		return false, err
	}
	return user, nil
}
