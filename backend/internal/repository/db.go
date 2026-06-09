package repository

import (
	"context"
	"fmt"

	"kasiraiai/backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("repository.InitDB: parse config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.DBMaxConnections)
	poolCfg.MinConns = int32(cfg.DBMaxIdleConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("repository.InitDB: connect: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("repository.InitDB: ping: %w", err)
	}

	return pool, nil
}
