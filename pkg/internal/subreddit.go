package internal

import (
	"time"

	"RD-Clone-API/pkg/model"
)

// NewSubreddit is a struct to create a new subreddit.
type NewSubreddit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SubredditResponse is the response to a created subreddit or getting a subreddit information.
type SubredditResponse struct {
	ID            int
	Name          string
	Description   string
	CreatedAt     time.Time
	UserID        int
	NumberOfPosts int
}

// BuildSubredditResponse constructs a SubredditResponse.
func BuildSubredditResponse(subR *model.Subreddit, postNum int) *SubredditResponse {
	return &SubredditResponse{
		ID:            subR.ID,
		Name:          subR.Name,
		Description:   subR.Description,
		CreatedAt:     subR.CreatedAt,
		UserID:        subR.UserID,
		NumberOfPosts: postNum,
	}
}
