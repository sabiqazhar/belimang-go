package service

import (
	"context"

	"github.com/sabiqazhar/belimang-go/user-service/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, request model.UserRegisterRequest, isAdmin bool) (string, error)
}
