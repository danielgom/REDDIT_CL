package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/routes/mock_services"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := &internal.RegisterResponse{
		Username: rr.Username,
		Email:    rr.Email,
		Enabled:  false,
	}

	mSvc := mock_services.NewMockUserService(ctrl)
	mSvc.EXPECT().SignUp(c.Request().Context(), &rr).Return(want, nil)

	h := NewUserHandler(mSvc)

	err = h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, err)

	responseString := rec.Body.String()

	require.Contains(t, responseString, want.Email)
	require.Contains(t, responseString, want.Username)
	require.Contains(t, responseString, fmt.Sprintf("%v", want.Enabled)) // There should be a better way to test this
}
