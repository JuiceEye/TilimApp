package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/dto/response"
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
	router.HandleFunc("POST /profile/xp", h.handleAddXP)
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

func (h *ProfileHandler) handleAddXP(w http.ResponseWriter, r *http.Request) {
	var req request.LessonCompletedRequest
	if err := utils.ParseBodyAndValidate(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	status, err := h.service.AddXPPoints(req.UserID, req.XPPoints)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	profile, status, err := h.service.GetProfile(req.UserID)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}
	resp := response.GetProfileResponse{
		UserID:           profile.UserID,
		Username:         profile.Username,
		RegistrationDate: profile.RegistrationDate,
		StreakDays:       profile.StreakDays,
		XPPoints:         profile.XPPoints,
		WordsLearned:     profile.WordsLearned,
		LessonsDone:      profile.LessonsDone,
	}
	utils.WriteJSONResponse(w, http.StatusOK, resp)
}
