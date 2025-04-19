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

type ModuleHandler struct {
	service *service.ModuleService
}

func NewModuleHandler(service *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		service: service,
	}
}

func (h *ModuleHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /main-page/module/{module_id}", h.handleGetMainPageModule)
}

func (h *ModuleHandler) handleGetMainPageModule(w http.ResponseWriter, r *http.Request) {
	payload := request.GetMainPageModuleRequest{}

	moduleID, err := strconv.ParseInt(r.PathValue("module_id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("module_id path param не найден"))
		return
	}
	payload.ModuleID = moduleID

	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	module, err := h.service.GetMainPageModuleByID(payload.ModuleID) // ?? или надо просто модуль и вызывать много сервисов??
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
	}

	moduleResponse := response.GetMainPageModuleResponse{
		ID:       module.ID,
		Title:    module.title,
		Sections: module.sections,
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, moduleResponse)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
