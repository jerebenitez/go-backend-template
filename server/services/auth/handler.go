package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/jerebenitez/go-backend-template/utils"
)

type HandlerFunc = func(http.ResponseWriter, *http.Request)

type Handler struct {
	Service  IAuthService
	Handlers map[string]HandlerFunc
}

type IAuthService interface {
	SignUp(User, string) (User, error)
}

func (h *Handler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	if err := utils.ParseJSON(r, &data); err != nil {
		slog.Error("error parsing body", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := User{Email: strings.ToLower(data.Email)}
	newUser, err := h.Service.SignUp(user, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.WriteJSON(w, 201, newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func NewAuthHandler(service IAuthService) *Handler {
	h := Handler{
		Service: service,
	}

	h.Handlers = map[string]HandlerFunc{
		"POST /signup": h.handleSignUp,
	}

	return &h
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	serviceRouter := http.NewServeMux()
	for route, handler := range h.Handlers {
		serviceRouter.HandleFunc(route, handler)
	}

	router.Handle("/auth/", http.StripPrefix("/auth", serviceRouter))
}
