package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	token2 "github.com/sabiqazhar/belimang-go/pkg/token"
	"github.com/sabiqazhar/belimang-go/user-service/internal/db"
	"github.com/sabiqazhar/belimang-go/user-service/internal/model"
	"github.com/sabiqazhar/belimang-go/user-service/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, request model.UserRegisterRequest, isAdmin bool) (string, error) {
	isValidEmail, err := s.userRepo.IsEmailAdminExists(ctx, request.Email)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error checking if email exists")
		return "", err
	}
	if isValidEmail {
		logger.Logger.Error().Msg("user with given email already exists")
		return "", errors.New("user with given email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error hashing password")
		return "", err
	}

	user := db.CreateUserParams{
		Email:    request.Email,
		Password: string(hashedPassword),
		Username: request.Username,
		IsAdmin:  pgtype.Bool{Bool: isAdmin},
	}

	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error creating user in repository")
		return "", err
	}

	token, err := token2.GenerateJWTToken(userID, "user")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error generating JWT token")
		return "", err
	}

	return token, nil
}
