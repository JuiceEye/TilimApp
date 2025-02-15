package handler

import (
	"net/http"
	"tilimauth/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /register/", h.Register)
	router.HandleFunc("POST /login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
