package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/infrasturcture/database"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	"github.com/sabiqazhar/belimang-go/user-service/internal/handler"
	"github.com/sabiqazhar/belimang-go/user-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/user-service/internal/service"
)

func main() {
	logger.InitLogger()

	r := gin.Default()

	dbConfig := database.PostgresConfig{
		Host:     "localhost",
		Port:     5433,
		User:     "user",
		Password: "password",
		DBName:   "user_service",
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	userRepo, err := repositories.NewUserRepository(db, 1)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to create user repository")
	}

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewHandler(r, userService)
	userHandler.RegisterRoutes()

	if err := r.Run(":8080"); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
