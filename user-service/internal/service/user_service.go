package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	token2 "github.com/sabiqazhar/belimang-go/pkg/token"
	errors2 "github.com/sabiqazhar/belimang-go/user-service/errors"
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
	existingUserByUsername, err := s.userRepo.GetUserByUsername(ctx, request.Username)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error checking if username exists")
		return "", err
	}

	if existingUserByUsername.ID != 0 {
		logger.Logger.Error().Msg("username already exists")
		return "", errors2.ErrEmailAlreadyExists
	}

	if isAdmin {
		existingAdmin, err := s.userRepo.IsEmailAdminExists(ctx, request.Email)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("error checking if admin email exists")
			return "", err
		}

		if existingAdmin {
			logger.Logger.Error().Msg("admin with given email already exists")
			return "", errors2.ErrAdminEmailExists
		}
	} else {
		existingUserByEmail, err := s.userRepo.GetUserByEmail(ctx, request.Email)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("error checking if email exists")
			return "", err
		}

		if existingUserByEmail.ID != 0 && !existingUserByEmail.IsAdmin.Bool {
			logger.Logger.Error().Msg("user with given email already exists")
			return "", errors2.ErrEmailAlreadyExists
		}
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
		IsAdmin:  pgtype.Bool{Bool: isAdmin, Valid: true},
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

func (s *UserServiceImpl) UserLogin(ctx context.Context, request model.UserLoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		logger.Logger.Error().Err(err).Msg("invalid password")
		return "", err
	}

	token, err := token2.GenerateJWTToken(user.ID, "user")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("error generating JWT token")
		return "", err
	}
	return token, nil
}
