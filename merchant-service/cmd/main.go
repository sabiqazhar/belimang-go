package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/infrasturcture/database"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/handler"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/service"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	pb "github.com/sabiqazhar/belimang-go/proto/merchant"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger.InitLogger()

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

	// --- gRPC Server Setup ---
	go func() {
		grpcPort := ":50051"
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %v", err)
		}

		s := grpc.NewServer()

		merchantGRPCHandler := handler.NewMerchantGRPCHandler(merchantService)
		pb.RegisterMerchantServiceServer(s, merchantGRPCHandler)

		logger.Logger.Info().Msgf("🚀 starting gRPC server on port %s", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := gin.Default()
	merchantHandler := handler.NewMerchantHandler(r, merchantService)
	merchantHandler.RegisterRoutes()

	httpPort := ":8081"
	logger.Logger.Info().Msgf("🚀 starting HTTP server on port %s", httpPort)
	if err := r.Run(httpPort); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to run HTTP server")
	}
}
