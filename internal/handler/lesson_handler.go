package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tilimauth/internal/dto/request"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
	"tilimauth/internal/utils"
	"time"
)

type LessonHandler struct {
	lessonService     *service.LessonService
	completionService *service.LessonCompletionService
}

func NewLessonHandler(service *service.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: service,
	}
}

func (h *LessonHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /lessons/{lesson_id}", h.handleGetLesson)
	router.HandleFunc("POST /lessons/{lesson_id}/complete", h.handleCompleteLesson)

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

	lesson, err := h.lessonService.GetLessonByID(payload.LessonID)
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

func (h *LessonHandler) handleCompleteLesson(w http.ResponseWriter, r *http.Request) {
	payload := &request.CompleteLessonRequest{}
	lessonIDPath := "lesson_id"
	lessonID, err := strconv.ParseInt(r.PathValue(lessonIDPath), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("%s path param не найден", lessonIDPath))
		return
	}
	payload.LessonID = lessonID
	payload.UserID = int64(32)

	if err := utils.ParseBodyAndValidate(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	lessonCompletion := &model.LessonCompletion{
		UserID:        payload.UserID,
		LessonID:      payload.LessonID,
		DateCompleted: time.Now(),
	}

	err = h.completionService.CompleteLesson(lessonCompletion)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = utils.WriteJSONResponse(w, http.StatusOK, map[string]int64{"lesson_id": lessonID})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
