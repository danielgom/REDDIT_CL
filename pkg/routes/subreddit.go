package routes

import (
	"net/http"

	"RD-Clone-API/pkg/context"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/labstack/echo/v4"
)

// SubRedditHandler is an instance of our reddit handler API.
type SubRedditHandler struct {
	redditSVC services.SubredditService
}

// NewSubRedditHandler returns a RedditHandler instance.
func NewSubRedditHandler(svc services.SubredditService) *SubRedditHandler {
	return &SubRedditHandler{redditSVC: svc}
}

// Register adds all routes related to user service.
func (h *SubRedditHandler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("/api/subreddit")
	authGroup.POST("", handler(h.Create))
	authGroup.GET("", handler(h.GetAll))
	authGroup.GET("/:id", handler(h.Get))
}

// Create creates a new subreddit and binds it to a user.
func (h *SubRedditHandler) Create(c context.Context) error {
	var req internal.NewSubreddit

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, redditErr := h.redditSVC.Create(c.Request().Context(), &req, c.User())
		if redditErr != nil {
			return nil, redditErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

// Get returns a single subreddit by ID.
func (h *SubRedditHandler) Get(c context.Context) error {
	var sRedditID int

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.redditSVC.Get(c.Request().Context(), sRedditID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

// GetAll returns all subreddits.
func (h *SubRedditHandler) GetAll(c context.Context) error {
	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getAllErr := h.redditSVC.GetAll(c.Request().Context())
		if getAllErr != nil {
			return nil, getAllErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}
