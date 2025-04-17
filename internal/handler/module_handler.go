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

type ModuleHandler struct {
	service *service.ModuleService
}

func NewModuleHandler(service *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		service: service,
	}
}

func (h *ModuleHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /module/{module_id}", h.getModule)
}

func (h *ModuleHandler) getModule(w http.ResponseWriter, r *http.Request) {
	payload := request.GetModuleRequest{}

	moduleID, err := strconv.ParseInt(r.PathValue("module_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user_id path param не найден"))
		return
	}
	payload.UserID = userID

	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	profile, status, err := h.service.GetProfile(payload.UserID)
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	profileResponse := response.GetModuleResponse{
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
