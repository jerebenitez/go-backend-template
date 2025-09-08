package users

import (
	"net/http"
)

type HandlerFunc = func(http.ResponseWriter, *http.Request)

type Handler struct {
	Handlers map[string]HandlerFunc
}

func NewUserHandler(service *UserService) *Handler {
	return &Handler{
		Handlers: service.RegisteredServices(),
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	serviceRouter := http.NewServeMux()
	for route, handler := range h.Handlers {
		serviceRouter.HandleFunc(route, handler)
	}

	router.Handle("/users/", http.StripPrefix("/users", serviceRouter))
}
