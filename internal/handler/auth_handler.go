package handler

import (
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
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload request.AuthLoginRequest
	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// todo: зашифровать пароль для сравнения но хз где это делается тут или в сервисе
	// err = utils.WriteJSONResponse(w, http.StatusOK, response)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusInternalServerError, err)
	// 	return
	// }
}

// todo: сделать логгирование стабильным (изучить log либо использовать только fmt, а не одно принта другое для ошибок
func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// todo: добавить валидацию для требований к паролю (спец. сиволы)
	payload := request.AuthRegistrationRequest{}
	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// todo: шифровать пароли
	user := model.User{
		Username:         payload.Username,
		Password:         payload.Password,
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

	token, err := auth.GenerateJWT(w, createdUser.ID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// todo: узнать как правильно возвращать токены: точно ли просто в body...?
	registerResponse := response.AuthRegistrationResponse{
		UserID: createdUser.ID,
		Token:  token,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, registerResponse)
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
