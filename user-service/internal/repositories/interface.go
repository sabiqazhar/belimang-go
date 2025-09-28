package repositories

import (
	"context"

	"github.com/sabiqazhar/belimang-go/user-service/internal/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user db.CreateUserParams) (int64, error)
	IsEmailAdminExists(ctx context.Context, email string) (bool, error)
}
