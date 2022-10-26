package services

import (
	"context"
	"time"

	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/samber/lo"
)

type postSvc struct {
	postDB       db.PostRepository
	userSvc      UserService
	subredditSvc SubredditService
}

// NewPostSvc creates a new Post service instance.
func NewPostSvc(postDB db.PostRepository, uSvc UserService, sSvc SubredditService) PostService {
	return &postSvc{postDB: postDB, userSvc: uSvc, subredditSvc: sSvc}
}

// Create creates a new post.
func (p *postSvc) Create(ctx context.Context, postReq *internal.NewPost, username string) (*internal.PostResponse, errors.CommonError) {
	user, getUserErr := p.userSvc.Get(ctx, username)
	if getUserErr != nil {
		return nil, getUserErr
	}

	subreddit, getSrErr := p.subredditSvc.Get(ctx, int(postReq.SubredditID))
	if getSrErr != nil {
		return nil, getSrErr
	}

	post := &model.Post{
		Name:        postReq.Title,
		URL:         "",
		Description: postReq.Description,
		VoteCount:   0,
		UserID:      user.ID,
		SubredditID: subreddit.ID,
		CreatedAt:   time.Now().Local(),
	}

	post, saveErr := p.postDB.Save(ctx, post)
	if saveErr != nil {
		return nil, saveErr
	}

	return internal.BuildPostResponse(post), nil
}

// Get gets a new post based on its ID.
func (p *postSvc) Get(ctx context.Context, postID int) (*internal.PostResponse, errors.CommonError) {
	post, getErr := p.postDB.FindByID(ctx, postID)
	if getErr != nil {
		return nil, getErr
	}
	return internal.BuildPostResponse(post), nil
}

// GetAllByUser retrieves all posts from a user.
func (p *postSvc) GetAllByUser(ctx context.Context, username string) ([]*internal.PostResponse, errors.CommonError) {
	user, getUserErr := p.userSvc.Get(ctx, username)
	if getUserErr != nil {
		return nil, getUserErr
	}

	return p.buildPostResponse(ctx, user.ID, p.postDB.FindAllByUser)
}

// GetAllBySubreddit retrieves all posts from a subreddit.
func (p *postSvc) GetAllBySubreddit(ctx context.Context, subredditID int) ([]*internal.PostResponse, errors.CommonError) {
	subreddit, getSubErr := p.subredditSvc.Get(ctx, subredditID)
	if getSubErr != nil {
		return nil, getSubErr
	}

	return p.buildPostResponse(ctx, subreddit.ID, p.postDB.FindAllBySubreddit)
}

// buildPostResponse retrieves all posts by query.
func (p *postSvc) buildPostResponse(ctx context.Context, query int, fn func(context.Context, int) ([]*model.Post, errors.CommonError)) ([]*internal.PostResponse, errors.CommonError) {
	posts, commonError := fn(ctx, query)
	if commonError != nil {
		return nil, commonError
	}

	return lo.Map(posts, func(post *model.Post, _ int) *internal.PostResponse {
		return internal.BuildPostResponse(post)
	}), nil
}
