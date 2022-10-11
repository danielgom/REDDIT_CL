// Package routes will be the responsible for adding all routes from the service.
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
	authGroup.POST("/refresh/token", handler(h.refreshToken))
}

// SignUp is used to create a new user.
func (h *UserHandler) SignUp(c context.Context) error {
	var req internal.RegisterRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, signErr := h.UsrSvc.SignUp(c.Request().Context(), &req)
		if signErr != nil {
			return nil, signErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

// VerifyAccount verifies an account based on the token that has been given.
func (h *UserHandler) VerifyAccount(c context.Context) error {
	token := c.Param("token")

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		verifyErr := h.UsrSvc.VerifyAccount(c.Request().Context(), token)
		if verifyErr != nil {
			return nil, verifyErr
		}
		return &context.GResponse{
			Status:   http.StatusOK,
			Response: map[string]interface{}{"account": "Validated", "status": http.StatusOK},
		}, nil
	})
}

// Login returns a JWT based on the user that has been logged in.
func (h *UserHandler) Login(c context.Context) error {
	var req internal.LoginRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, logErr := h.UsrSvc.Login(c.Request().Context(), &req)
		if logErr != nil {
			return nil, logErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

func (h *UserHandler) refreshToken(c context.Context) error {
	var req internal.RefreshTokenRequest

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, refreshErr := h.UsrSvc.RefreshToken(c.Request().Context(), &req)
		if refreshErr != nil {
			return nil, refreshErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}
