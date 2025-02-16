package handler

import (
	"log"
	"net/http"
	"strconv"
	"tilimauth/internal/auth"
	"tilimauth/internal/dto"
	"tilimauth/internal/model"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
	"time"
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
	router.HandleFunc("POST /register", h.handleRegister)
	router.HandleFunc("POST /login", h.handleLogin)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//todo: пофиксить маршалинг постман запроса и узнать как проверять целостност body (нарушая синтаксис постмана летит хуйня)
	var payload dto.AuthRegistrationRequest
	if err := utils.ParseJSONRequest(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Printf("Parsing request %v", payload)
	//todo: добавить валидацию для требований к паролю
	if err := payload.Validate(); err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}
	//todo: кодировка паролей
	user := model.User{
		Username:         payload.Username,
		Password:         payload.Password,
		Email:            payload.Email,
		PhoneNumber:      payload.PhoneNumber,
		Image:            payload.Image,
		RegistrationDate: time.Now(),
	}

	createdUser, err, status := h.service.Register(user)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	token, err := auth.GenerateJWT(w, createdUser.Id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	response := dto.AuthRegistrationResponse{
		Id:    createdUser.Id,
		Token: token,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"user_id": strconv.Itoa(response.Id),
		"token":   response.Token,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *AuthHandler) handleProtectedRoute(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.VerifyJWT(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Access allowed",
		"user_id": strconv.Itoa(userID),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
