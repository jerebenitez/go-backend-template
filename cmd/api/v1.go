package api

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	addr string
	db *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	subrouter := http.NewServeMux()

	// Register routes for each service here:
	// authHandler := auth.NewHandler()
	// authHandler.RegisterRoutes(subrouter)

	// This can be handled by nginx, so it might not be needed
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", subrouter))

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
