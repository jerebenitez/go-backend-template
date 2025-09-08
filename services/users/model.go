package users

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	Id         string `json:"id"`
	UserName   string `json:"username"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  pgtype.Timestamp `json:"createdAt"`
	UpdatedAt  pgtype.Timestamp `json:"updatedAt"`
}
