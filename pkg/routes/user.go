package routes

import (
	"RD-Clone-API/pkg/internal"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserService interface {
	SignUp(context.Context, *internal.RegisterRequest) error
}

type UserHandler struct {
	UsrSvc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{UsrSvc: svc}
}

func (h *UserHandler) Register(r *echo.Echo) {
	r.POST("/signup", h.SignUp)
}

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
