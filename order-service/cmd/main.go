package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/infrasturcture/database"
	"github.com/sabiqazhar/belimang-go/order-service/internal/handler"
	"github.com/sabiqazhar/belimang-go/order-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/order-service/internal/service"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

func main() {
	logger.InitLogger()

	r := gin.Default()

	orderNode, err := snowflake.NewNode(3)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to create snowflake node for merchant service")
	}

	dbConfig := database.PostgresConfig{
		Host:     "localhost",
		Port:     5433,
		User:     "user",
		Password: "password",
		DBName:   "order_service",
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	//dataStore := repositories.NewStore(db)
	orderRepo := repositories.NewOrderRepository(db, orderNode)
	orderService := service.NewOrderServiceImpl(orderRepo, db)
	orderHandler := handler.NewHandler(r, orderService)
	orderHandler.RegisterRoutes()

	logger.Logger.Info().Msg("starting order service on port 8082")
	if err := r.Run(":8083"); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to run server")
	}
}
