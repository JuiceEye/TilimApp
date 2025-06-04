package handler

import (
	"errors"
	"net/http"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /subscriptions/purchase", h.handleSubscriptionPurchase)
}

func (h *SubscriptionHandler) handleSubscriptionPurchase(w http.ResponseWriter, r *http.Request) {
	payload := &request.SubscriptionPurchaseRequest{}

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := utils.GetUserID(r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	subscriptionID, err = h.service.BuySubscription(userID, payload.ExpiresAt)
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

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]int64{"subscription_id": subscriptionID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
