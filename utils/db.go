package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	User 	 string
	Password string
	Path 	 string
	Name 	 string
	DSN		 string
}

func NewPool(cfg DbConfig) (*pgxpool.Pool, *context.Context, error) {
	var ctx = context.Background()
	pool, err := pgxpool.New(ctx, getConnectionString(cfg))
	if err != nil {
		return nil, nil, err	
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, nil, err
	}

	return pool, &ctx, nil
}

func getConnectionString(cfg DbConfig) string {
	if cfg.DSN == "" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s/%s",
			cfg.User,
			cfg.Password,
			cfg.Path,
			cfg.Name,
		)
	} else {
		return cfg.DSN
	}
}
