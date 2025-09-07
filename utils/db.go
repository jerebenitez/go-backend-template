package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	User string
	Password string
	Path string
	Name string
}

var ctx = context.Background()

func NewPool(cfg DbConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, getConnectionString(cfg))
	if err != nil {
		return nil, err	
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func getConnectionString(cfg DbConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Path,
		cfg.Name,
	)
}
