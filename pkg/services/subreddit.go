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

type subredditSvc struct {
	subRDB  db.SubredditRepository
	userSvc UserService
}

// NewSubredditSvc creates a new Subreddit service instance.
func NewSubredditSvc(subRDB db.SubredditRepository, svc UserService) SubredditService {
	return &subredditSvc{subRDB: subRDB, userSvc: svc}
}

// Create creates a new subreddit and binds a user to it as a creator.
func (s *subredditSvc) Create(ctx context.Context, subreddit *internal.NewSubreddit, username string) (*internal.SubredditResponse,
	errors.CommonError) {
	currentTime := time.Now().Local()

	user, getUsrErr := s.userSvc.Get(ctx, username)
	if getUsrErr != nil {
		return nil, getUsrErr
	}

	subR := &model.Subreddit{
		Name:        subreddit.Name,
		Description: subreddit.Description,
		CreatedAt:   currentTime,
		UserID:      user.ID,
	}

	savedSubR, saveErr := s.subRDB.Save(ctx, subR)
	if saveErr != nil {
		return nil, saveErr
	}

	return internal.BuildSubredditResponse(savedSubR), nil
}

// Get gets the subreddit.
func (s *subredditSvc) Get(ctx context.Context, subRedditID int) (*internal.SubredditResponse, errors.CommonError) {
	subRedditResponse, commonError := s.subRDB.FindByID(ctx, subRedditID)
	if commonError != nil {
		return nil, commonError
	}

	return internal.BuildSubredditResponse(subRedditResponse), nil
}

// GetAll returns a list of all subreddits.
func (s *subredditSvc) GetAll(ctx context.Context) ([]*internal.SubredditResponse, errors.CommonError) {
	subRList, commonError := s.subRDB.FindAll(ctx)
	if commonError != nil {
		return nil, commonError
	}

	return lo.Map(subRList, func(subRR *model.Subreddit, _ int) *internal.SubredditResponse {
		return internal.BuildSubredditResponse(subRR)
	}), nil
}
