package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/handler"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/upload-service/internal/service"
)

func main() {
	// DB
	conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost:5433/upload_service")
	if err != nil {
		log.Fatal("DB:", err)
	}
	defer conn.Close(context.Background())

	repo, err := repositories.NewUploadRepository(conn, 1)
	if err != nil {
		log.Fatal("Repo:", err)
	}

	svc, err := service.NewUploadService(
		repo,
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		"images",
		"http://localhost:9000",
		false,
	)
	if err != nil {
		log.Fatal("Service:", err)
	}

	r := gin.Default()
	// TODO: add auth middleware that sets c.Set("userID", id)
	r.POST("/image", handler.NewUploadHandler(svc).UploadImage)

	log.Println("Upload service running on :8080")
	r.Run(":8083")
}