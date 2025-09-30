package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(config PostgresConfig) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Logger.Info().Str("host", config.Host).Int("port", config.Port).Msg("Connected to PostgreSQL database successfully")
	return conn, nil
}
