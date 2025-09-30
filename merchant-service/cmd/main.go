package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/infrasturcture/database"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/handler"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/service"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

func main() {
	logger.InitLogger()

	r := gin.Default()

	merchantNode, err := snowflake.NewNode(2)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to create snowflake node for merchant service")
	}

	dbConfig := database.PostgresConfig{
		Host:     "localhost",
		Port:     5433,
		User:     "user",
		Password: "password",
		DBName:   "merchant_service",
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	merchantRepo := repositories.NewMerchantRepository(db, merchantNode)
	merchantService := service.NewMerchantService(merchantRepo)
	merchantHandler := handler.NewMerchantHandler(r, merchantService)

	merchantHandler.RegisterRoutes()

	logger.Logger.Info().Msg("starting merchant service on port 8081")
	if err := r.Run(":8081"); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to run server")
	}
}
