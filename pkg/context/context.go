// Package context is responsible for the addition powers to echo.Context.
//
//nolint:wrapcheck // Should not wrap echo JSON error
package context

import (
	"fmt"
	"net/http"

	"RD-Clone-API/pkg/util/errors"
	"github.com/labstack/echo/v4"
)

// Context is a custom echo context.
type Context struct {
	echo.Context
}

// BindAndValidate binds and validates structs if required.
func (c *Context) BindAndValidate(req any, fn func() error) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	fmt.Println("here omg")
	err := fn()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

// Handler turns a echo.HandlerFunc into a custom handler of ours.
func Handler(fn func(Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return fn(Context{Context: c})
	}
}
