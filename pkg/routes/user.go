// Package routes will be the responsible for adding all routes from the service.
package routes

import (
	"RD-Clone-API/pkg/services"
	"errors"
	"net/http"

	"RD-Clone-API/pkg/internal"

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
		return c.JSON(http.StatusBadRequest, errors.New("bad request"))
	}

	err := req.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.UsrSvc.SignUp(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, "user has been registered")
}
