package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/auth"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/dto/response"
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
	router.HandleFunc("POST /refresh", h.handleRefreshToken) // это для обновления токенов эндпоинт
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload request.AuthLoginRequest
	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, status, err := h.service.Login(payload.Username, payload.Password)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	tokens, err := auth.GenerateTokenPair(w, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	loginResponse := response.AuthLoginResponse{
		UserID: user.ID,
		Tokens: tokens,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, loginResponse)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

// todo: сделать логгирование стабильным (изучить log либо использовать только fmt, а не одно принта другое для ошибок
func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {

	// todo: добавить валидацию для требований к паролю (спец. сиволы)
	payload := request.AuthRegistrationRequest{}
	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err))
		return
	}

	user := model.User{
		Username:         payload.Username,
		Password:         hashedPassword,
		Email:            payload.Email,
		PhoneNumber:      payload.PhoneNumber,
		Image:            payload.Image,
		RegistrationDate: time.Now(),
	}

	createdUser, status, err := h.service.Register(user)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	tokens, err := auth.GenerateTokenPair(w, createdUser.ID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// todo: узнать как правильно возвращать токены: точно ли просто в body...?
	registerResponse := response.AuthRegistrationResponse{
		UserID: createdUser.ID,
		Tokens: tokens,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, registerResponse)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *AuthHandler) handleProtectedRoute(w http.ResponseWriter, r *http.Request) {
	createdUser, err := auth.VerifyTokens(r, "access")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Access allowed",
		"user_id": strconv.FormatInt(createdUser, 10),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

// только для получения аксесс токена, поэтому не хендлим защищённый маршрут как для аксесса
func (h *AuthHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.VerifyTokens(r, "refresh")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	tokens, err := auth.GenerateTokenPair(w, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refreshResponse := response.AuthRegistrationResponse{
		UserID: userID,
		Tokens: tokens,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, refreshResponse)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}
