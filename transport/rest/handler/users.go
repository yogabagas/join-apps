package handler

import (
	"encoding/json"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/shared/util"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"net/http"
	"strconv"
)

// CreateUsers handler
// @Summary Create New User
// @Description New User Registration
// @Tags Users
// @Produce json
// @Param users body service.CreateUsersReq true "Request Create User"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/users [POST]
func (h *HandlerImpl) CreateUsers(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodPost {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.CreateUsersReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	err := h.Controller.UsersController.CreateUsers(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)
}

// Logout handler
// @Summary Lgout
// @Description Logout endpoint
// @Tags Users
// @Produce json
//
// @Success 200 {object} response.JSONResponse().APIStatusSuccess()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/logout [POST]
func (h *HandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	res := response.NewJSONResponse()

	if r.Method != http.MethodDelete {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	token := r.Header.Get("Authorization")

	claims, err := util.GetUserData(token)
	if err != nil {
		res.SetError(response.ErrUnauthorized).SetMessage(err.Error()).Send(w)
		return
	}

	req := service.LogoutReq{
		UserUID: claims.UserUID,
	}

	err = h.Controller.UsersController.Logout(r.Context(), req)

	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusNoContent().Send(w)
}

// GetUsersWithPagination handler
// @Summary GetUsersWithPagination
// @Description GetUsersWithPagination for get users detail with limit
// @Tags Users
// @Produce json
// @Param name query string false "user fullname e.g John Doe"
// @Param limit query int false "limit data; default 10"
// @Param page query int false "number of page; default 1"
// @Success 200 {object} response.JSONResponse{data=service.GetUsersWithPaginationResp}
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/users [GET]
func (h *HandlerImpl) GetUsersWithPagination(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodGet {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.GetUsersWithPaginationReq

	if name := r.URL.Query().Get("name"); name != "" {
		req.Fullname = name
	}

	var limitToInt int
	if limit := r.URL.Query().Get("limit"); limit != "" {
		limitToInt, _ = strconv.Atoi(limit)
	}

	if limitToInt <= 0 {
		limitToInt = 10
	}
	req.Limit = limitToInt

	var pageToInt int
	if page := r.URL.Query().Get("page"); page != "" {
		pageToInt, _ = strconv.Atoi(page)
	}

	if pageToInt <= 0 {
		pageToInt = 1
	}
	req.Page = pageToInt

	resp, err := h.Controller.UsersController.GetUsersWithPagination(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	res.SetData(resp).Send(w)
}
