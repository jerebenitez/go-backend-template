package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
	ctx  *context.Context
}

func NewUserRepository(pool *pgxpool.Pool, ctx *context.Context) *UserRepository {
	return &UserRepository{
		pool: pool,
		ctx: ctx,
	}
}

func (r UserRepository) GetAllUsers() ([]User, error) {
	sql := `SELECT * FROM users`

	rows, err := r.pool.Query(*r.ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.UserName,
			&user.Email,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
