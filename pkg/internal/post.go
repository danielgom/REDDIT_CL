package internal

import (
	"time"

	"RD-Clone-API/pkg/model"
)

// NewPost is a struct to create a Post.
type NewPost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SubredditID uint
}

// PostResponse contains the Post information.
type PostResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	VoteCount   int       `json:"vote_count"`
	UserID      int       `json:"user_id"`
	CreatedDate time.Time `json:"created_date"`
	SubredditID int       `json:"subreddit_id"`
}

// BuildPostResponse constructs a PostResponse from a post model.
func BuildPostResponse(post *model.Post) *PostResponse {
	return &PostResponse{
		ID:          post.ID,
		Title:       post.Name,
		URL:         post.URL,
		Description: post.Description,
		VoteCount:   post.VoteCount,
		UserID:      post.UserID,
		CreatedDate: post.CreatedAt,
		SubredditID: post.SubredditID,
	}
}
