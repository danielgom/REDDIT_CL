// Package routes will be the responsible for adding all routes from the service.
//
//nolint:wrapcheck // Should not wrap echo JSON error
package routes

import (
	"net/http"

	"RD-Clone-API/pkg/context"
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
func (h *UserHandler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("/api/auth")
	authGroup.POST("/signup", handler(h.SignUp))
	authGroup.GET("/accountVerification/:token", handler(h.VerifyAccount))
	authGroup.POST("/login", handler(h.Login))
}

// SignUp is used to create a new user.
func (h *UserHandler) SignUp(c context.Context) error {
	var req internal.RegisterRequest

	return c.BindAndValidate(&req, func() error {
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
		}

		response, signErr := h.UsrSvc.SignUp(c.Request().Context(), &req)
		if signErr != nil {
			return c.JSON(signErr.Status(), signErr)
		}

		return c.JSON(http.StatusCreated, response)
	})
}

// VerifyAccount verifies an account based on the token that has been given.
func (h *UserHandler) VerifyAccount(c context.Context) error {
	token := c.Param("token")

	verifyErr := h.UsrSvc.VerifyAccount(c.Request().Context(), token)
	if verifyErr != nil {
		return c.JSON(http.StatusInternalServerError, verifyErr)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"account": "Validated", "status": http.StatusOK})
}

// Login returns a JWT based on the user that has been logged in.
func (h *UserHandler) Login(c context.Context) error {
	var req internal.LoginRequest

	return c.BindAndValidate(&req, func() error {
		response, logErr := h.UsrSvc.Login(c.Request().Context(), &req)
		if logErr != nil {
			return c.JSON(logErr.Status(), logErr)
		}

		return c.JSON(http.StatusCreated, response)
	})
}
