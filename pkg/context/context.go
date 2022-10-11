// Package context is responsible for the addition powers to echo.Context.
package context

import (
	"net/http"

	"RD-Clone-API/pkg/util/errors"
	"github.com/labstack/echo/v4"
)

// Context is a custom echo context.
type Context struct {
	echo.Context
}

// GResponse is a wrapper for the response and status of controllers.
type GResponse struct {
	Status   int
	Response any
}

// BindAndValidateResp binds and validates structs if required.
func (c *Context) BindAndValidateResp(req any, fn func() (*GResponse, errors.CommonError)) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid json format", err))
	}

	res, err := fn()

	return c.parseResponse(res, err)
}

// NoBindResp use this when no json body is being provided/required.
func (c *Context) NoBindResp(fn func() (*GResponse, errors.CommonError)) error {
	res, err := fn()

	return c.parseResponse(res, err)
}

func (c *Context) parseResponse(res *GResponse, err errors.CommonError) error {
	if err != nil {
		return c.JSON(err.Status(), err)
	}

	return c.JSON(res.Status, res.Response)
}

// Handler turns an echo.HandlerFunc into a custom handler of ours.
func Handler(fn func(Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return fn(Context{Context: c})
	}
}
