package utils

import (
	"context"
	"fmt"

	"github.com/thd-spatial-ai/ignis/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildConnectionString(cfg config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
}

// ConnectPool creates a connection pool to the database
func ConnectPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
