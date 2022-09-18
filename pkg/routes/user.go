// Package routes will be the responsible for adding all routes from the service.
//
//nolint:wrapcheck // Should not wrap echo JSON error
package routes

import (
	"net/http"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/labstack/echo/v4"
)

// UserHandler is an instance of our user handler API.
type UserHandler struct {
	UsrSvc services.UserService
}

// NewUserHandler returns a UserHandler instance.
func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{UsrSvc: svc}
}

// Register adds all routes related to user service.
func (h *UserHandler) Register(r *echo.Echo) {
	r.POST("/signup", h.SignUp)
}

// SignUp is used to create a new user.
func (h *UserHandler) SignUp(c echo.Context) error {
	var req internal.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	err := req.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	response, signErr := h.UsrSvc.SignUp(c.Request().Context(), &req)
	if signErr != nil {
		return c.JSON(signErr.Status(), signErr)
	}

	return c.JSON(http.StatusCreated, response)
}
