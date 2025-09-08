package users

import (
	"fmt"
	"net/http"

	"github.com/jerebenitez/go-backend-template/utils"
)

type IUserRepository interface {
	GetAllUsers() ([]User, error)
	CreateNewUser(User) (User, error)
	DeleteUser(string) error
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
		"GET /": s.GetUsers,
		"POST /": s.CreateUser,
		"DELETE /{id}": s.DeleteUser,
	}
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := utils.ParseJSON(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := s.repo.CreateNewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.repo.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, fmt.Sprintf("User %s deleted!", id)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
