package users

import (
	"fmt"
	"net/http"

	"github.com/jerebenitez/go-backend-template/utils"
)

type HandlerFunc = func(http.ResponseWriter, *http.Request)

type Handler struct {
	Service IUserService
	Handlers map[string]HandlerFunc
}

type IUserService interface {
	CreateUser(User) (User, error)
	GetUsers() ([]User, error)
	DeleteUser(string) error
}

func (h *Handler) handleDeleteUsers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	 if err := h.Service.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, fmt.Sprintf("User %s deleted!", id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.Service.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := utils.ParseJSON(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 200, newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func NewUserHandler(service IUserService) *Handler {
	h := Handler{
		Service: service,
	}

	h.Handlers = map[string]HandlerFunc{
		"POST /": h.handleCreateUser,
		"GET /": h.handleGetUsers,
		"DELETE /{id}": h.handleDeleteUsers,
	}

	return &h
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	serviceRouter := http.NewServeMux()
	for route, handler := range h.Handlers {
		serviceRouter.HandleFunc(route, handler)
	}

	router.Handle("/users/", http.StripPrefix("/users", serviceRouter))
}
