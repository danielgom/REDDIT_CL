package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"RD-Clone-API/pkg/context"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/labstack/echo/v4"
)

var (
	errInvalidID = fmt.Errorf("invalid id")
)

// PostHandler is an instance of our reddit handler API.
type PostHandler struct {
	postSvc services.PostService
}

// NewPostHandler returns a PostHandler instance.
func NewPostHandler(svc services.PostService) *PostHandler {
	return &PostHandler{postSvc: svc}
}

// Register adds all routes related to user service.
func (h *PostHandler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("/api/post")
	authGroup.POST("", handler(h.Create))
	authGroup.GET("/user", handler(h.GetAllByUser))
	authGroup.GET("/subreddit/:id", handler(h.GetAllBySubreddit))
	authGroup.GET("/:id", handler(h.Get))
}

// Create saves a new Post to a subreddit.
func (h *PostHandler) Create(c context.Context) error {
	var req internal.NewPost

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, createErr := h.postSvc.Create(c.Request().Context(), &req, c.User())
		if createErr != nil {
			return nil, createErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

// Get retrieves a post by ID.
func (h *PostHandler) Get(c context.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", errInvalidID)
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.postSvc.Get(c.Request().Context(), postID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// GetAllByUser gets all posts users by user.
func (h *PostHandler) GetAllByUser(c context.Context) error {
	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.postSvc.GetAllByUser(c.Request().Context(), c.User())
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// GetAllBySubreddit gets all posts by subreddit id.
func (h *PostHandler) GetAllBySubreddit(c context.Context) error {
	subredditID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", errInvalidID)
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.postSvc.GetAllBySubreddit(c.Request().Context(), subredditID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}
