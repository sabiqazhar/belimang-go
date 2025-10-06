package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func NewPostgresDB(config PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	// 1. Parse the DSN string into a config object
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	// 2. Set pool parameters (adjust these values based on your needs)
	poolConfig.MaxConns = 25                      // Max number of connections in the pool
	poolConfig.MinConns = 5                       // Min number of connections to keep open
	poolConfig.MaxConnLifetime = time.Hour        // Max lifetime of a connection
	poolConfig.MaxConnIdleTime = time.Minute * 30 // Time a connection can be idle

	// 3. Create the pool using the modified config
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// ... ping and return the pool as before ...
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Logger.Info().Str("host", config.Host).Int("port", config.Port).Msg("Connected to PostgreSQL database successfully")
	return pool, nil
}
