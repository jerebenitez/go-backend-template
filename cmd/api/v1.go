package api

import (
	"database/sql"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
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
