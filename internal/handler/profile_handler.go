package handler

//
// import (
// 	"net/http"
// 	"strconv"
// 	"tilimauth/internal/auth"
// 	"tilimauth/internal/dto/request"
// 	"tilimauth/internal/dto/response"
// 	"tilimauth/internal/service"
// 	"tilimauth/internal/utils"
// )
//
// type ProfileHandler struct {
// 	service *service.ProfileService
// }
//
// func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
// 	return &ProfileHandler{
// 		service: service,
// 	}
// }
//
// func (h *ProfileHandler) RegisterRoutes(router *http.ServeMux) {
// 	router.HandleFunc("POST /profile/{id}", h.handleReadProfile)
// }
//
// func (h *ProfileHandler) handleReadProfile(w http.ResponseWriter, r *http.Request) {
// 	payload := request.ReadProfileRequest{}
//
// 	if err := parseReadProfilePathParams(r, &payload); err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}
//
// 	if err := utils.ParseBodyAndValidate(r, &payload); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}
//
// 	readProfile, status, err := h.service.ReadProfile(userID)
// 	if err != nil {
// 		utils.WriteError(w, status, err)
// 		return
// 	}
//
// 	token, err := auth.GenerateJWT(w, createdUser.ID)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusUnauthorized, err)
// 		return
// 	}
//
// 	// todo: узнать как правильно возвращать токены: точно ли просто в body...?
// 	response := response.AuthRegistrationResponse{
// 		UserID: createdUser.ID,
// 		Token:  token,
// 	}
//
// 	err = utils.WriteJSONResponse(w, http.StatusOK, response)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}
// }
//
// func parseReadProfilePathParams(r *http.Request, payload *request.ReadProfileRequest) error {
// 	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	payload.UserID = userID
//
// 	return nil
// }
