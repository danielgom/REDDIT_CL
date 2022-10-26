package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/context"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/routes/mock_services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

var (
	errCreateSubreddit = fmt.Errorf("could not create subreddit")
	errGetSubreddit    = fmt.Errorf("could not get subreddit")
	errGetSubreddits   = fmt.Errorf("could not query subreddits")
)

func TestSubredditHandler(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, h *SubRedditHandler, e *echo.Echo, srSvc *mock_services.MockSubredditService){
		"test create subreddit success":           testCreateSubreddit,
		"test create subreddit svc error fails":   testCreateSubredditSvcErr,
		"test get subreddit success":              testGetSubreddit,
		"test get subreddit svc error fails":      testGetSubredditSvcErr,
		"test get all subreddits success":         testGetAllSR,
		"test get all subreddits svc error fails": testGetAllSRSvcErr,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			handler, ec, svc, teardown := setupSubredditSvc(t)
			defer teardown()
			fn(t, handler, ec, svc)
		})
	}
}

func setupSubredditSvc(t *testing.T) (*SubRedditHandler, *echo.Echo, *mock_services.MockSubredditService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	userSvcMock := mock_services.NewMockSubredditService(ctrl)

	handler := NewSubRedditHandler(userSvcMock)
	e := echo.New()
	v := config.GetValidator()
	err := config.AddValidators(v.Validator)
	require.NoError(t, err)
	e.Validator = v

	return handler, e, userSvcMock, func() {
		defer ctrl.Finish()
	}
}

func testCreateSubreddit(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	tSubreddit := internal.NewSubreddit{
		Name:        "testsubr",
		Description: "this is going to be my test subr",
	}

	userJSON, err := json.Marshal(tSubreddit)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/subreddit", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	want := &internal.SubredditResponse{
		ID:            1,
		Name:          "testsubr",
		Description:   "this is going to be my test subr",
		UserID:        2,
		NumberOfPosts: 2,
	}

	svc.EXPECT().Create(c.Request().Context(), &tSubreddit, gomock.Any()).Return(want, nil)

	err = h.Create(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, err)

	var got internal.SubredditResponse

	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Description, got.Description)
	require.Equal(t, want.NumberOfPosts, got.NumberOfPosts)
}

func testCreateSubredditSvcErr(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	tSubreddit := internal.NewSubreddit{
		Name:        "testsubr",
		Description: "this is going to be my test subr",
	}

	userJSON, err := json.Marshal(tSubreddit)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/subreddit", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().Create(c.Request().Context(), &tSubreddit, gomock.Any()).Return(nil,
		errors.NewBadRequestError("test err", errCreateSubreddit))

	err = h.Create(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "could not create subreddit")
}

func testGetSubreddit(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/subreddit", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	want := &internal.SubredditResponse{
		ID:            10,
		Name:          "testgetsubr",
		Description:   "amazing subreddit!",
		UserID:        5,
		NumberOfPosts: 21,
	}

	svc.EXPECT().Get(c.Request().Context(), gomock.Any()).Return(want, nil)

	err := h.Get(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, err)

	var got internal.SubredditResponse

	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, want.Name, got.Name)
	require.Equal(t, want.ID, got.ID)
	require.Equal(t, want.Description, got.Description)
	require.Equal(t, want.NumberOfPosts, got.NumberOfPosts)
}

func testGetSubredditSvcErr(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/subreddit/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().Get(c.Request().Context(), gomock.Any()).Return(nil,
		errors.NewBadRequestError("test err", errGetSubreddit))

	err := h.Get(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "could not get subreddit")
}

func testGetAllSR(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/subreddit", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	want := []*internal.SubredditResponse{
		{
			ID:            10,
			Name:          "testgetsubr",
			Description:   "amazing subreddit!",
			UserID:        5,
			NumberOfPosts: 21,
		},
		{
			ID:            11,
			Name:          "testgetsubr2",
			Description:   "amazing subreddit2!",
			UserID:        5,
			NumberOfPosts: 50,
		},
		{
			ID:            12,
			Name:          "testgetsubr3",
			Description:   "amazing subreddit3!",
			UserID:        6,
			NumberOfPosts: 100,
		},
	}

	svc.EXPECT().GetAll(c.Request().Context()).Return(want, nil)

	err := h.GetAll(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var got []*internal.SubredditResponse

	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func testGetAllSRSvcErr(t *testing.T, h *SubRedditHandler, e *echo.Echo, svc *mock_services.MockSubredditService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/subreddit", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().GetAll(c.Request().Context()).Return(nil,
		errors.NewBadRequestError("test err", errGetSubreddits))

	err := h.GetAll(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "could not query subreddits")
}
