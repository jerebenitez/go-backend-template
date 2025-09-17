package api

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jerebenitez/go-backend-template/services/auth"
	"github.com/jerebenitez/go-backend-template/utils"
)

type APIServer struct {
	addr string
	pool *pgxpool.Pool
	ctx  *context.Context
}

func NewAPIServer(addr string, pool *pgxpool.Pool, ctx *context.Context) *APIServer {
	return &APIServer{
		addr: addr,
		pool: pool,
		ctx:  ctx,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	// Register routes for each service here:
	authRepo := auth.NewAuthRepository(s.pool, s.ctx)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)
	authHandler.RegisterRoutes(router)

	// This can be handled by nginx
	//router.Handle("/api/v1/", http.StripPrefix("/api/v1", subrouter))

	// Health check endpoint
	router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.WriteJSON(w, http.StatusOK, "OK"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
