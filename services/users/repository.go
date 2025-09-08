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

func (r *UserRepository) GetAllUsers() ([]User, error) {
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

func (r *UserRepository) CreateNewUser(user User) (User, error) {
	sql := `
		INSERT INTO users(username, email)
		VALUES ($1, $2)
		RETURNING *
	`

	var newUser User
	err := r.pool.QueryRow(*r.ctx, sql, user.UserName, user.Email).Scan(
		&newUser.Id,
		&newUser.UserName,
		&newUser.Email,
		&newUser.IsVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("error creating user: %w", err)
	}

	return newUser, nil
}

func (r *UserRepository) DeleteUser(id string) error {
	sql := `
		DELETE FROM users WHERE id = $1
	`

	_, err := r.pool.Query(*r.ctx, sql, id)
	if err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}
