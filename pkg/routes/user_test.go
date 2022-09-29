package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/routes/mock_services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

var errSign = fmt.Errorf("could not sign up")
var errVerifyDB = fmt.Errorf("could not verify account")

func TestUserHandler(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, h *UserHandler, uSvc *mock_services.MockUserService){
		"test user sing up success":             testCreateUser,
		"test create user bad request fails":    testCreateUserBadJSON,
		"test create user invalid fields fails": testCreateUserValidation,
		"test create user service err fails ":   testCreateUserSvcErr,
		"test verify account success":           testVerifyAccount,
		"test verify account service err fails": testVerifyAccountSvcErr,
		"test login success":                    testLogin,
		"test login service err fails":          testLoginSvcErr,
		"test login bad request fails":          testLoginBadJSON,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			handler, svc, teardown := setupUserSvc(t)
			defer teardown()
			fn(t, handler, svc)
		})
	}
}

func setupUserSvc(t *testing.T) (*UserHandler, *mock_services.MockUserService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	userSvcMock := mock_services.NewMockUserService(ctrl)

	handler := NewUserHandler(userSvcMock)

	return handler, userSvcMock, func() {
		defer ctrl.Finish()
	}
}

func testCreateUser(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	rr := internal.RegisterRequest{
		Username: "Daniel",
		Password: "P@ssw0rd123123",
		Email:    "dga_355@hotmail.com",
	}

	userJSON, err := json.Marshal(rr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	want := &internal.RegisterResponse{
		Username: rr.Username,
		Email:    rr.Email,
		Enabled:  false,
	}

	svc.EXPECT().SignUp(c.Request().Context(), &rr).Return(want, nil)

	err = h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, err)

	responseString := rec.Body.String()

	require.Contains(t, responseString, want.Email)
	require.Contains(t, responseString, want.Username)
	require.Contains(t, responseString, fmt.Sprintf("%v", want.Enabled)) // There should be a better way to test this
}

func testCreateUserBadJSON(t *testing.T, h *UserHandler, _ *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader("123123{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SignUp(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "invalid json format")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func testCreateUserValidation(t *testing.T, h *UserHandler, _ *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SignUp(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "invalid fields")
}

func testCreateUserSvcErr(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	rr := internal.RegisterRequest{
		Username: "Daniel",
		Password: "P@ssw0rd123123",
		Email:    "dga_355@hotmail.com",
	}

	userJSON, err := json.Marshal(rr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc.EXPECT().SignUp(c.Request().Context(), &rr).Return(nil, errors.NewBadRequestError("test err", errSign))

	err = h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "could not sign up")
}

func testVerifyAccount(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/accountVerification/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc.EXPECT().VerifyAccount(c.Request().Context(), gomock.Any()).Return(nil)

	err := h.VerifyAccount(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Validated")
}

func testVerifyAccountSvcErr(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/accountVerification/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc.EXPECT().VerifyAccount(c.Request().Context(), gomock.Any()).
		Return(errors.NewInternalServerError("database error", errVerifyDB))

	err := h.VerifyAccount(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "database error")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func testLogin(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	lr := internal.LoginRequest{
		UserOrEmail: "dga_355@hotmail.com",
		Password:    "Test12345@",
	}

	userJSON, err := json.Marshal(lr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	want := &internal.LoginResponse{
		Username:     "DanielGA",
		Email:        lr.UserOrEmail,
		Token:        "jsontoken123",
		RefreshToken: "jsonrefresh123",
		ExpiresAt:    time.Now(),
	}

	svc.EXPECT().Login(c.Request().Context(), &lr).Return(want, nil)

	err = h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, err)

	responseString := rec.Body.String()

	require.Contains(t, responseString, want.Email)
	require.Contains(t, responseString, want.Username)
	require.Contains(t, responseString, want.Token)
	require.Contains(t, responseString, want.RefreshToken)
}

func testLoginSvcErr(t *testing.T, h *UserHandler, svc *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	lr := internal.LoginRequest{
		UserOrEmail: "dga_355@hotmail.com",
		Password:    "Test12345@",
	}

	userJSON, err := json.Marshal(lr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc.EXPECT().Login(c.Request().Context(), &lr).Return(nil, errors.NewUnauthorisedError("service error"))

	err = h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "service error")
}

func testLoginBadJSON(t *testing.T, h *UserHandler, _ *mock_services.MockUserService) {
	t.Helper()
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("123123{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid json format")
}
