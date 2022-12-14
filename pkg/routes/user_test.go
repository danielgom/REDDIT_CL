package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"RD-Clone-API/pkg/config"
	"RD-Clone-API/pkg/context"
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
	for scenario, fn := range map[string]func(t *testing.T, h *UserHandler, e *echo.Echo, uSvc *mock_services.MockUserService){
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
			handler, ec, svc, teardown := setupUserSvc(t)
			defer teardown()
			fn(t, handler, ec, svc)
		})
	}
}

func setupUserSvc(t *testing.T) (*UserHandler, *echo.Echo, *mock_services.MockUserService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	userSvcMock := mock_services.NewMockUserService(ctrl)

	handler := NewUserHandler(userSvcMock)
	e := echo.New()
	v := config.GetValidator()
	err := config.AddValidators(v.Validator)
	require.NoError(t, err)
	e.Validator = v

	return handler, e, userSvcMock, func() {
		defer ctrl.Finish()
	}
}

func testCreateUser(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

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
	c := context.Context{Context: e.NewContext(req, rec)}

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

	var got internal.RegisterResponse

	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, want.Email, got.Email)
	require.Equal(t, want.Username, got.Username)
	require.Equal(t, want.Enabled, got.Enabled)
}

func testCreateUserBadJSON(t *testing.T, h *UserHandler, e *echo.Echo, _ *mock_services.MockUserService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader("123123{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	err := h.SignUp(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "invalid json format")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func testCreateUserValidation(t *testing.T, h *UserHandler, e *echo.Echo, _ *mock_services.MockUserService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader("{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	err := h.SignUp(c)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "Field validation")
}

func testCreateUserSvcErr(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

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
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().SignUp(c.Request().Context(), &rr).Return(nil, errors.NewBadRequestError("test err", errSign))

	err = h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "could not sign up")
}

func testVerifyAccount(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/accountVerification/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().VerifyAccount(c.Request().Context(), gomock.Any()).Return(nil)

	err := h.VerifyAccount(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Validated")
}

func testVerifyAccountSvcErr(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/accountVerification/10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().VerifyAccount(c.Request().Context(), gomock.Any()).
		Return(errors.NewInternalServerError("database error", errVerifyDB))

	err := h.VerifyAccount(c)
	require.NoError(t, err)
	require.Contains(t, rec.Body.String(), "database error")
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func testLogin(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

	lr := internal.LoginRequest{
		UserOrEmail: "dga_355@hotmail.com",
		Password:    "Test12345@",
	}

	userJSON, err := json.Marshal(lr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	want := &internal.LoginResponse{
		Username:     "DanielGA",
		Email:        lr.UserOrEmail,
		Token:        "jsontoken123",
		RefreshToken: "jsonrefresh123",
		ExpiresAt:    time.Now(),
	}

	svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(want, nil)

	err = h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, err)

	var got internal.LoginResponse
	err = json.NewDecoder(rec.Body).Decode(&got)
	require.NoError(t, err)

	require.Contains(t, want.Email, got.Email)
	require.Contains(t, want.Username, got.Username)
	require.Contains(t, want.Token, got.Token)
	require.Contains(t, want.RefreshToken, got.RefreshToken)
}

func testLoginSvcErr(t *testing.T, h *UserHandler, e *echo.Echo, svc *mock_services.MockUserService) {
	t.Helper()

	lr := internal.LoginRequest{
		UserOrEmail: "dga_355@hotmail.com",
		Password:    "Test12345@",
	}

	userJSON, err := json.Marshal(lr)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.NewUnauthorisedError("service error"))

	err = h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "service error")
}

func testLoginBadJSON(t *testing.T, h *UserHandler, e *echo.Echo, _ *mock_services.MockUserService) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("123123{}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := context.Context{Context: e.NewContext(req, rec)}

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid json format")
}
