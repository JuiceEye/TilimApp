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

type LessonHandler struct {
	service *service.ProfileService
}

func NewLessonHandler(service *service.LessonService) *LessonHandler {
	return &LessonHandler{
		service: service,
	}
}

func (h *LessonHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /profile/{user_id}", h.handleGetLesson)
}

func (h *ProfileHandler) handleGetLesson(w http.ResponseWriter, r *http.Request) {
	payload := &request.GetLessonRequest{}
	lessonIDPath := "lesson_id"
	lessonID, err := strconv.ParseInt(r.PathValue(lessonIDPath), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s path param не найден", lessonIDPath))
		return
	}
	payload.lessonID = lessonID

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	lesson, status, err := h.service.GetLessonByID(payload.lessonID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, lesson)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
