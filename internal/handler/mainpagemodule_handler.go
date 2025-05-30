package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type ModuleHandler struct {
	service *service.MainPageModuleService
}

func NewMainPageModuleHandler(service *service.MainPageModuleService) *ModuleHandler {
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

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	module, err := h.service.GetMainPageModuleByID(payload.ModuleID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, module)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
