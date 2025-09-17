package auth

import (
	"fmt"
	"net/mail"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id              string             `json:"id,omitempty"`
	Email           string             `json:"email,omitempty"`
	EmailVerified   bool               `json:"emailVerified,omitempty"`
	EmailVerifiedAt pgtype.Timestamptz `json:"emailVerifiedAt,omitempty"`
	PasswordHash    string             `json:"passwordHash,omitempty"`
	IsActive        bool               `json:"isActive,omitempty"`
	IsSuperuser     bool               `json:"isSuperuser,omitempty"`
	LastLogin       pgtype.Timestamptz `json:"lastLogin,omitempty"`
	MfaEnabled      bool               `json:"mfaEnabled,omitempty"`
	PhoneNumber     string             `json:"phoneNumber,omitempty"`
	PhoneVerifiedAt pgtype.Timestamptz `json:"phoneVerifiedAt,omitempty"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt,omitempty"`
	UpdatedAt       pgtype.Timestamptz `json:"updatedAt,omitempty"`
}

func ValidateUser(u User) error {
	_, err := mail.ParseAddress(u.Email)
	return err
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must contain at least 8 characters")
	}

	containsUpper := false
	containsLower := false
	containsSymbol := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			containsUpper = true
		} else if unicode.IsLower(r) {
			containsLower = true
		} else if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			containsSymbol = true
		}
	}

	if !containsLower {
		return fmt.Errorf("password must contain at least one lower-case character")
	}

	if !containsUpper {
		return fmt.Errorf("password must contain at least one upper-case character")
	}

	if !containsSymbol {
		return fmt.Errorf("password must contain at least one symbol")
	}

	return nil
}
