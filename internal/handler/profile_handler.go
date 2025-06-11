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
	router.HandleFunc("GET /me", h.handleGetMyProfile)
	router.HandleFunc("PATCH /profile/avatar", h.handleUpdateProfilePicture)
	router.HandleFunc("PATCH /profile/username", h.handleUpdateUsername)
	router.HandleFunc("PATCH /profile/email", h.handleUpdateEmail)
	router.HandleFunc("PATCH /profile/password", h.handleUpdatePassword)
	router.HandleFunc("GET /profile/activity/{user_id}", h.handleGetUserActivity)
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	payload := &request.GetProfileRequest{}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user_id path param не найден"))
		return
	}
	payload.UserID = userID

	if err = utils.ParseBodyAndValidate(r, payload); err != nil {
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
		Image:            profile.Image,
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

func (h *ProfileHandler) handleGetMyProfile(w http.ResponseWriter, r *http.Request) {
	payload := &request.GetProfileRequest{}

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	payload.UserID = userID

	if err = utils.ParseBodyAndValidate(r, payload); err != nil {
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
		Image:            profile.Image,
		StreakDays:       profile.StreakDays,
		XPPoints:         profile.XPPoints,
		WordsLearned:     profile.WordsLearned,
		LessonsDone:      profile.LessonsDone,
		IsSubscribed:     &profile.IsSubscribed,
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

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.UpdateProfilePicture(userID, payload.Image)
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

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.UpdateUsername(userID, payload.Username)
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

func (h *ProfileHandler) handleUpdateEmail(w http.ResponseWriter, r *http.Request) {
	payload := &request.UpdateEmailRequest{}

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.UpdateEmail(userID, payload.Email, payload.Password)
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

func (h *ProfileHandler) handleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	payload := &request.UpdatePasswordRequest{}

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.UpdatePassword(userID, payload.Password, payload.NewPassword)
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

	if err := utils.WriteJSONResponse(w, http.StatusOK, map[string]int64{"user_id": userID}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *ProfileHandler) handleGetUserActivity(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user_id path param не найден"))
		return
	}

	if userID == 58 {
		// Hardcoded dataset for user 58
		activity := []response.UserActivityResponse{
			{Date: "2025-03-01", LessonsCompleted: 15},
			{Date: "2025-03-02", LessonsCompleted: 3},
			{Date: "2025-03-03", LessonsCompleted: 10},
			{Date: "2025-03-04", LessonsCompleted: 2},
			{Date: "2025-03-05", LessonsCompleted: 4},
			{Date: "2025-03-06", LessonsCompleted: 1},
			{Date: "2025-03-07", LessonsCompleted: 1},
			{Date: "2025-03-08", LessonsCompleted: 0},
			{Date: "2025-03-09", LessonsCompleted: 3},
			{Date: "2025-03-10", LessonsCompleted: 3},
			{Date: "2025-03-11", LessonsCompleted: 1},
			{Date: "2025-03-12", LessonsCompleted: 5},
			{Date: "2025-03-13", LessonsCompleted: 1},
			{Date: "2025-03-14", LessonsCompleted: 4},
			{Date: "2025-03-15", LessonsCompleted: 3},
			{Date: "2025-03-16", LessonsCompleted: 5},
			{Date: "2025-03-17", LessonsCompleted: 5},
			{Date: "2025-03-18", LessonsCompleted: 5},
			{Date: "2025-03-19", LessonsCompleted: 0},
			{Date: "2025-03-20", LessonsCompleted: 11},
			{Date: "2025-03-21", LessonsCompleted: 4},
			{Date: "2025-03-22", LessonsCompleted: 3},
			{Date: "2025-03-23", LessonsCompleted: 14},
			{Date: "2025-03-24", LessonsCompleted: 5},
			{Date: "2025-03-25", LessonsCompleted: 4},
			{Date: "2025-03-26", LessonsCompleted: 3},
			{Date: "2025-03-27", LessonsCompleted: 1},
			{Date: "2025-03-28", LessonsCompleted: 3},
			{Date: "2025-03-29", LessonsCompleted: 4},
			{Date: "2025-03-30", LessonsCompleted: 4},
			{Date: "2025-03-31", LessonsCompleted: 2},
			{Date: "2025-04-01", LessonsCompleted: 0},
			{Date: "2025-04-02", LessonsCompleted: 3},
			{Date: "2025-04-03", LessonsCompleted: 0},
			{Date: "2025-04-04", LessonsCompleted: 1},
			{Date: "2025-04-05", LessonsCompleted: 3},
			{Date: "2025-04-06", LessonsCompleted: 0},
			{Date: "2025-04-07", LessonsCompleted: 2},
			{Date: "2025-04-08", LessonsCompleted: 1},
			{Date: "2025-04-09", LessonsCompleted: 4},
			{Date: "2025-04-10", LessonsCompleted: 0},
			{Date: "2025-04-11", LessonsCompleted: 0},
			{Date: "2025-04-12", LessonsCompleted: 1},
			{Date: "2025-04-13", LessonsCompleted: 2},
			{Date: "2025-04-14", LessonsCompleted: 5},
			{Date: "2025-04-15", LessonsCompleted: 2},
			{Date: "2025-04-16", LessonsCompleted: 3},
			{Date: "2025-04-17", LessonsCompleted: 3},
			{Date: "2025-04-18", LessonsCompleted: 4},
			{Date: "2025-04-19", LessonsCompleted: 3},
			{Date: "2025-04-20", LessonsCompleted: 2},
			{Date: "2025-04-21", LessonsCompleted: 3},
			{Date: "2025-04-22", LessonsCompleted: 1},
			{Date: "2025-04-23", LessonsCompleted: 4},
			{Date: "2025-04-24", LessonsCompleted: 1},
		}

		err := utils.WriteJSONResponse(w, http.StatusOK, activity)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	activity, err := h.service.GetUserActivity(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, activity)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
