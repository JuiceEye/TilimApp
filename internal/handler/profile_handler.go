package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/dto/response"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type ProfileHandler struct {
	service *service.ProfileService
}

func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

func (h *ProfileHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /profile/{user_id}", h.handleGetProfile)
	router.HandleFunc("PUT /profile/avatar", h.handleUpdateProfilePicture)
	router.HandleFunc("PUT /profile/username", h.handleUpdateUsername)
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	payload := &request.GetProfileRequest{}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user_id path param не найден"))
		return
	}
	payload.UserID = userID

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	profile, status, err := h.service.GetProfile(payload.UserID)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	profileResponse := response.GetProfileResponse{
		UserID:           profile.UserID,
		Username:         profile.Username,
		RegistrationDate: profile.RegistrationDate,
		StreakDays:       profile.StreakDays,
		XPPoints:         profile.XPPoints,
		WordsLearned:     profile.WordsLearned,
		LessonsDone:      profile.LessonsDone,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, profileResponse)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *ProfileHandler) handleUpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	payload := &request.UpdateProfilePictureRequest{}

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := int64(32)

	err := h.service.UpdateProfilePicture(userID, payload.Image)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]int64{"user_id": userID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *ProfileHandler) handleUpdateUsername(w http.ResponseWriter, r *http.Request) {
	payload := &request.UpdateUsernameRequest{}

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := int64(32)

	err := h.service.UpdateUsername(userID, payload.Username)
	if err != nil {
		var bre *service.BadRequestError
		if errors.As(err, &bre) {
			utils.WriteError(w, http.StatusBadRequest, err)
		} else if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]int64{"user_id": userID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
