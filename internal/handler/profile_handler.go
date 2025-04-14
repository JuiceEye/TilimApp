package handler

import (
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
	router.HandleFunc("POST /profile/{id}", h.handleReadProfile)
}

func (h *ProfileHandler) handleReadProfile(w http.ResponseWriter, r *http.Request) {
	payload := request.ReadProfileRequest{}

	if err := parseReadProfilePathParams(r, &payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	profile, status, err := h.service.GetProfile(payload.UserID)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	profileResponse := response.ReadProfileResponse{
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

func parseReadProfilePathParams(r *http.Request, payload *request.ReadProfileRequest) error {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		return err
	}
	payload.UserID = userID

	return nil
}
