package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"tilimauth/internal/auth"
	"tilimauth/internal/dto"
	"tilimauth/internal/dto/request"
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
	if err := utils.ParseJSONRequest(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if !(payload.Username == "test" && payload.Password == "qwerty") &&
		!(payload.Username == "JuiceEye" && payload.Password == "qwerty") {
		err := utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "wrong username or password"})
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		return
	}

	rand.Seed(time.Now().UnixNano()) // Seed to ensure randomness
	randomNumber := rand.Intn(100)   // Generates a random number between 0 and 99
	token, err := auth.GenerateJWT(w, randomNumber)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	response := dto.AuthLoginResponse{
		Token: token,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

// todo: пофиксить маршалинг постман запроса и узнать как проверять целостность body
// (нарушая синтаксис постмана json.decode пытается раскодировать кривой джейсон и ответ ошибки выходит некрасивый)
// todo: сделать логгирование стабильным (изучить log либо использовать только fmt, а не одно принта другое для ошибок
// todo: сделать уникальное логгирование не привязанное к handle методу а к любому запросу
func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//todo: добавить валидацию для требований к паролю (спец. сиволы)
	var payload = request.AuthRegistrationRequest{}
	if err := utils.ParseAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := payload.Validate(); err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	//todo: шифровать пароли
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

	//todo: узнать как правильно возвращать токены: точно ли просто в body...?
	response := dto.AuthRegistrationResponse{
		UserId: createdUser.Id,
		Token:  token,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, response)
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
