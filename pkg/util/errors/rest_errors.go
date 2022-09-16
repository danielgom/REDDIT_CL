// Package errors contains a common way to declare API errors
package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errInvalidJSON = errors.New("invalid json")
)

// RestErr is a common interface where we can create different kind of errors.
type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type commonError struct {
	M string        `json:"message"`
	S int           `json:"status"`
	E string        `json:"error"`
	C []interface{} `json:"causes"`
}

func (r *commonError) Message() string {
	return r.M
}

func (r *commonError) Status() int {
	return r.S
}

func (r *commonError) Error() string {
	return r.E
}

func (r *commonError) Causes() []interface{} {
	return r.C
}

// NewRestError creates an error based on user input.
func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return &commonError{
		M: message,
		S: status,
		E: err,
		C: causes,
	}
}

// NewRestErrorFromBytes is based on upcoming bytes.
func NewRestErrorFromBytes(b []byte) (RestErr, error) {
	var r commonError
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errInvalidJSON
	}
	return &r, nil
}

// NewBadRequestError returns a bad request RestErr.
func NewBadRequestError(m string) RestErr {
	return &commonError{
		M: m,
		S: http.StatusBadRequest,
		E: "Bad Request",
	}
}

// NewNotFoundError returns a bad request RestErr
// We normally error not found errors based on database response.
func NewNotFoundError(m string) RestErr {
	return &commonError{
		M: m,
		S: http.StatusNotFound,
		E: "Not found",
	}
}

// NewInternalServerError returns an internal server error (such as DB calls or service-to-service calls).
func NewInternalServerError(m string, err error) RestErr {
	r := &commonError{
		M: m,
		S: http.StatusInternalServerError,
		E: "Internal server error",
	}
	if err != nil {
		r.C = append(r.C, err.Error())
	}
	return r
}

// NewUnauthorisedError returns an Unauthorised error (no access to the api).
func NewUnauthorisedError(m string) RestErr {
	return &commonError{
		M: m,
		S: http.StatusUnauthorized,
		E: "Unauthorised",
	}
}
