package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
	ctx  *context.Context
}

func NewAuthRepository(pool *pgxpool.Pool, ctx *context.Context) *AuthRepository {
	return &AuthRepository{
		pool: pool,
		ctx:  ctx,
	}
}

func (r *AuthRepository) CreateNewUser(user User) (User, error) {
	sql := `
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, created_at
	`

	var newUser User
	row := r.pool.QueryRow(*r.ctx, sql, user.Email, user.PasswordHash)
	err := row.Scan(
		&newUser.Id,
		&newUser.Email,
		&newUser.CreatedAt,
	)
	if err != nil {
		slog.Error("cannot create user", "error", err)
		// TODO: Check error, return a more user-friendly version
		return User{}, fmt.Errorf("error creating user: %w", err)
	}

	return newUser, nil
}
