package routes

import (
	"RD-Clone-API/pkg/routes/mock_services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"RD-Clone-API/pkg/internal"

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

	mSvc := mock_services.NewMockUserService(ctrl)
	mSvc.EXPECT().SignUp(c.Request().Context(), &rr).Return(nil)

	h := NewUserHandler(mSvc)

	err = h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "\"user has been registered\"", strings.TrimSpace(rec.Body.String()))
}
