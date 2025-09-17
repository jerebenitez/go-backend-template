package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST = 14

type IAuthRepository interface {
	CreateNewUser(User) (User, error)
}

type AuthService struct {
	repo IAuthRepository
}

func NewAuthService(repo IAuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) SignUp(user User, password string) (User, error) {
	if err := ValidatePassword(password); err != nil {
		return User{}, err
	}

	if err := ValidateUser(user); err != nil {
		return User{}, err
	}

	if len([]byte(password)) > 72 {
		return User{}, fmt.Errorf("password cannot be longer than 72 bytes")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return User{}, err
	}
	user.PasswordHash = string(hashedPassword)

	newUser, err := s.repo.CreateNewUser(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}
