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

type LessonHandler struct {
	service *service.LessonService
}

func NewLessonHandler(service *service.LessonService) *LessonHandler {
	return &LessonHandler{
		service: service,
	}
}

func (h *LessonHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /lessons/{lesson_id}", h.handleGetLesson)
}

func (h *LessonHandler) handleGetLesson(w http.ResponseWriter, r *http.Request) {
	payload := &request.GetLessonRequest{}
	lessonIDPath := "lesson_id"
	lessonID, err := strconv.ParseInt(r.PathValue(lessonIDPath), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s path param не найден", lessonIDPath))
		return
	}
	payload.LessonID = lessonID

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	lesson, err := h.service.GetLessonByID(payload.LessonID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, lesson)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
