package users

import (
	"net/http"

	"github.com/jerebenitez/go-backend-template/utils"
)

type IUserRepository interface {
	GetAllUsers() ([]User, error)
}

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
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
	if r.Method == http.MethodGet {
		users, err := s.repo.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.WriteJSON(w, 200, users); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.NotFound(w, r)
}
