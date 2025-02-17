package api

import (
	"database/sql"
	"log"
	"net/http"
	"tilimauth/internal/handler"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
)

type Server struct {
	address string
	db      *sql.DB
}

func NewServer(address string, db *sql.DB) *Server {
	return &Server{
		address: address,
		db:      db,
	}
}

func (s *Server) Run() error {
	router := http.NewServeMux()

	userRepo := repository.NewAuthRepo(s.db)
	userService := service.NewAuthService(userRepo)
	userHandler := handler.NewAuthHandler(userService)
	userHandler.RegisterRoutes(router)

	log.Printf("Starting server on %s...", s.address)

	return http.ListenAndServe(s.address, router)
}
