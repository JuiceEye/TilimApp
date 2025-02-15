package handler

import (
	"net/http"
	"tilimauth/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /register", h.Register)
	router.HandleFunc("POST /login", h.Login)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	 
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}
