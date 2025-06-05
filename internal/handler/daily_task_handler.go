package handler

import (
	"net/http"
	"tilimauth/internal/dto/response"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type DailyTaskHandler struct {
	service *service.DailyTaskService
}

func NewDailyTaskHandler(service *service.DailyTaskService) *DailyTaskHandler {
	return &DailyTaskHandler{
		service: service,
	}
}

func (h *DailyTaskHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /daily-tasks", h.handleGetDailyTasks)
}

// handleGetDailyTasks handles the request to get daily tasks for the current user
func (h *DailyTaskHandler) handleGetDailyTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tasks, err := h.service.GetUserDailyTasks(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := response.ToDailyTaskResponseList(tasks)
	if err := utils.WriteJSONResponse(w, http.StatusOK, response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}