package users

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (r *UserRepository) GetAllUsers() {}
