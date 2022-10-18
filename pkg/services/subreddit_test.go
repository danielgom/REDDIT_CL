package services

import (
	"context"
	"fmt"
	"testing"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/services/mock_repositories"
	"RD-Clone-API/pkg/services/mock_services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const testUsername = "testingUname"

var (
	errSaveSubReddit         = fmt.Errorf("could not save subreddit")
	errGetSubReddit          = fmt.Errorf("could not get subreddit")
	errGetSubRedditPostCount = fmt.Errorf("could not get post count")
	errFindSubreddit         = fmt.Errorf("could not get subreddit list")
)

func TestRedditService(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, srRepo *mock_repositories.MockSubredditRepository,
		userSvc *mock_services.MockUserService, svc SubredditService){
		"test create subreddit success":           testCreateReddit,
		"test create subreddit no user fails":     testCreateSubredditNoUser,
		"test create subreddit save error fails":  testCreateSubredditSaveErr,
		"test get subreddit success":              testGetSubreddit,
		"test get subreddit get err fails":        testGetSubredditErr,
		"test get subreddit get count err fails":  testGetSubredditCountErr,
		"test get all subreddits success":         testGetAllSubreddits,
		"test get all subreddits find err fails":  testGetAllSubredditsFindErr,
		"test get all subreddits count err fails": testGetAllSubredditsCountErr,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			subRRepo, uSvc, svc, teardown := setupRedditSvc(t)
			defer teardown()
			fn(t, subRRepo, uSvc, svc)
		})
	}
}

func setupRedditSvc(t *testing.T) (*mock_repositories.MockSubredditRepository, *mock_services.MockUserService,
	SubredditService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	subredditRepo := mock_repositories.NewMockSubredditRepository(ctrl)
	userService := mock_services.NewMockUserService(ctrl)

	svc := NewSubredditSvc(subredditRepo, userService)

	return subredditRepo, userService, svc, func() {
		svc = nil
		defer ctrl.Finish()
	}
}

func testCreateReddit(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, userSvc *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	uWant := &internal.UserResponse{
		ID: 1,
	}

	subWant := &model.Subreddit{
		ID:          1,
		Name:        "test1",
		Description: "Some awesome description",
		UserID:      1,
	}

	userSvc.EXPECT().Get(ctx, testUsername).Return(uWant, nil)
	srRepo.EXPECT().Save(ctx, gomock.Any()).Return(subWant, nil)

	req := &internal.NewSubreddit{
		Name:        "test1",
		Description: "Some awesome description",
	}

	res, createErr := svc.Create(ctx, req, "testingUname")
	require.NoError(t, createErr)
	require.NotNil(t, res)
	require.Equal(t, "Some awesome description", res.Description)
	require.Equal(t, 1, res.ID)
	require.Equal(t, 1, res.UserID)
	require.Equal(t, 0, res.NumberOfPosts)
}

func testCreateSubredditNoUser(t *testing.T, _ *mock_repositories.MockSubredditRepository, userSvc *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	userSvc.EXPECT().Get(ctx, testUsername).Return(nil, errors.NewNotFoundError("user not found"))

	req := &internal.NewSubreddit{
		Name:        "test1",
		Description: "Some awesome description",
	}

	res, createErr := svc.Create(ctx, req, "testingUname")
	require.Nil(t, res)
	require.NotNil(t, createErr)
	require.Equal(t, "user not found", createErr.Message())
}

func testCreateSubredditSaveErr(t *testing.T, srRepo *mock_repositories.MockSubredditRepository,
	userSvc *mock_services.MockUserService, svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	uWant := &internal.UserResponse{
		ID: 1,
	}

	userSvc.EXPECT().Get(ctx, testUsername).Return(uWant, nil)
	srRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil, errors.NewInternalServerError("Database error", errSaveSubReddit))

	req := &internal.NewSubreddit{
		Name:        "test1",
		Description: "Some awesome description",
	}

	res, createErr := svc.Create(ctx, req, "testingUname")
	require.Nil(t, res)
	require.NotNil(t, createErr)
}

func testGetSubreddit(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	want := &model.Subreddit{
		ID:          1,
		Name:        "amazing subreddit",
		Description: "cool description",
		UserID:      1,
	}

	srRepo.EXPECT().FindByID(ctx, 1).Return(want, nil)
	srRepo.EXPECT().GetSubredditPostCount(ctx, want.ID).Return(10, nil)

	res, createErr := svc.Get(ctx, 1)
	require.NoError(t, createErr)
	require.NotNil(t, res)
	require.Equal(t, want.ID, res.ID)
	require.Equal(t, want.Name, res.Name)
	require.Equal(t, want.Description, res.Description)
	require.Equal(t, want.UserID, res.UserID)
	require.Equal(t, 10, res.NumberOfPosts)
}

func testGetSubredditErr(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()
	ctx := context.Background()

	srRepo.EXPECT().FindByID(ctx, 1).Return(nil, errors.NewInternalServerError("Database error", errGetSubReddit))

	res, createErr := svc.Get(ctx, 1)
	require.NotNil(t, createErr)
	require.Nil(t, res)
}

func testGetSubredditCountErr(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	want := &model.Subreddit{
		ID:          1,
		Name:        "amazing subreddit",
		Description: "cool description",
		UserID:      1,
	}

	srRepo.EXPECT().FindByID(ctx, 1).Return(want, nil)
	srRepo.EXPECT().GetSubredditPostCount(ctx, want.ID).Return(0, errors.NewInternalServerError("Database error",
		errGetSubRedditPostCount))

	res, createErr := svc.Get(ctx, 1)
	require.Nil(t, res)
	require.NotNil(t, createErr)
}

func testGetAllSubreddits(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	wantList := []*model.Subreddit{{
		ID:          1,
		Name:        "amazing subreddit",
		Description: "cool description",
		UserID:      1,
	}, {
		ID:          2,
		Name:        "amazing subreddit2",
		Description: "cool description2",
		UserID:      2,
	}, {
		ID:          3,
		Name:        "amazing subreddit2",
		Description: "cool description2",
		UserID:      3,
	}, {
		ID:          4,
		Name:        "amazing subreddit2",
		Description: "cool description2",
		UserID:      4,
	}}

	srRepo.EXPECT().FindAll(ctx).Return(wantList, nil)
	srRepo.EXPECT().GetSubredditPostCount(gomock.Any(), 1).Return(10, nil)
	srRepo.EXPECT().GetSubredditPostCount(gomock.Any(), 2).Return(15, nil)
	srRepo.EXPECT().GetSubredditPostCount(gomock.Any(), 3).Return(20, nil)
	srRepo.EXPECT().GetSubredditPostCount(gomock.Any(), 4).Return(25, nil)

	res, createErr := svc.GetAll(ctx)
	require.NoError(t, createErr)
	require.Equal(t, 4, len(res))
}

func testGetAllSubredditsFindErr(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	srRepo.EXPECT().FindAll(ctx).Return(nil, errors.NewInternalServerError("Database error",
		errFindSubreddit))

	res, createErr := svc.GetAll(ctx)
	require.NotNil(t, createErr)
	require.Nil(t, res)
}

func testGetAllSubredditsCountErr(t *testing.T, srRepo *mock_repositories.MockSubredditRepository, _ *mock_services.MockUserService,
	svc SubredditService) {
	t.Helper()

	ctx := context.Background()

	wantList := []*model.Subreddit{{
		ID:          1,
		Name:        "amazing subreddit",
		Description: "cool description",
		UserID:      1,
	}}

	srRepo.EXPECT().FindAll(ctx).Return(wantList, nil)
	srRepo.EXPECT().GetSubredditPostCount(gomock.Any(), 1).Return(0, errors.NewInternalServerError("Database error",
		errGetSubRedditPostCount))

	res, createErr := svc.GetAll(ctx)
	require.NotNil(t, createErr)
	require.Nil(t, res)
}
