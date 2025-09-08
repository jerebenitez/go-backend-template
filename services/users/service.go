package users

import (
	"net/http"
)

type IUserRepository interface {
	GetAllUsers() []string
}

type UserService struct {
	repo *IUserRepository
}

func NewUserService(repo *IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) RegisteredServices() map[string]HandlerFunc {
	return map[string]HandlerFunc{
		"/": s.GetUsers,
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	
}
