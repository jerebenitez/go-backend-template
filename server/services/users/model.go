package users

import (
	"net/mail"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id         string `db:"id" json:"id,omitempty"`
	UserName   string `db:"username" json:"username"`
	Email      string `db:"email" json:"email"`
	IsVerified bool   `db:"is_verified" json:"isVerified,omitempty"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at" json:"updatedAt,omitempty"`
}

func ValidateUser(u User) error {
	_, err := mail.ParseAddress(u.Email)
	return err
}
