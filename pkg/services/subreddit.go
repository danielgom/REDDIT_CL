package services

import (
	"context"
	"time"

	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util/errors"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
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

	return internal.BuildSubredditResponse(savedSubR, 0), nil
}

// Get gets the subreddit.
func (s *subredditSvc) Get(ctx context.Context, subRedditID int) (*internal.SubredditResponse, errors.CommonError) {
	subRedditResponse, commonError := s.subRDB.FindByID(ctx, subRedditID)
	if commonError != nil {
		return nil, commonError
	}

	count, countErr := s.subRDB.GetSubredditPostCount(ctx, subRedditResponse.ID)
	if countErr != nil {
		return nil, countErr
	}

	return internal.BuildSubredditResponse(subRedditResponse, count), nil
}

// GetAll returns a list of all subreddits.
func (s *subredditSvc) GetAll(ctx context.Context) ([]*internal.SubredditResponse, errors.CommonError) {
	subRList, commonError := s.subRDB.FindAll(ctx)
	if commonError != nil {
		return nil, commonError
	}

	g, gCtx := errgroup.WithContext(ctx)

	resChan := make(chan *internal.SubredditResponse, len(subRList))

	// Probably needs channels in order to work properly (append to slice is not really concurrent safe)
	for _, subreddit := range subRList {
		func(ctx context.Context, sr *model.Subreddit) {
			g.Go(func() error {
				count, e := s.subRDB.GetSubredditPostCount(ctx, sr.ID)
				if e != nil {
					return e
				}
				resChan <- internal.BuildSubredditResponse(sr, count)
				return nil
			})
		}(gCtx, subreddit)
	}

	err := g.Wait()
	if err != nil {
		return nil, errors.NewInternalServerError("could not retrieve subreddits", err)
	}

	close(resChan)

	return lo.ChannelToSlice(resChan), nil
}
