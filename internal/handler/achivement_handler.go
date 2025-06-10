package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/achievement"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type AchievementHandler struct {
	service *achievement.AchievementService
}

func NewAchievementHandler(service *achievement.AchievementService) *AchievementHandler {
	return &AchievementHandler{
		service: service,
	}
}

func (h *AchievementHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /achievements/{user_id}", h.handleGetAchievements)
}

func (h *AchievementHandler) handleGetAchievements(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user_id path param не найден"))
		return
	}

	achievements, err := h.service.GetAchievements(userID)
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

	if err := utils.WriteJSONResponse(w, http.StatusOK, achievements); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
