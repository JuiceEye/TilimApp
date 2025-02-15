package handler

import (
	"net/http"
	"tilimauth/internal/dto"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
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
	router.HandleFunc("POST /register/", h.handleRegister)
	router.HandleFunc("POST /login", h.handleLogin)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload dto.AuthRegistrationRequest
	if err := utils.ParseJSONRequest(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

}
