package handler

import (
	"net/http"
	"tilimauth/internal/dto/response"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type LeaderboardsHandler struct {
	service *service.LeaderboardsService
}

func NewLeaderboardsHandler(service *service.LeaderboardsService) *LeaderboardsHandler {
	return &LeaderboardsHandler{
		service: service,
	}
}

func (h *LeaderboardsHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /leaderboards", h.handleGetLeaderboards)
}

func (h *LeaderboardsHandler) handleGetLeaderboards(w http.ResponseWriter, r *http.Request) {
	leaderboards, status, err := h.service.GetLeaderboards()
	if err != nil {
		utils.WriteError(w, status, err)
		return
	}

	responseList := response.ToProfileResponseList(leaderboards)
	err = utils.WriteJSONResponse(w, http.StatusOK, responseList)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
