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
