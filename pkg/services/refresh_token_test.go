package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/services/mock_repositories"
	"RD-Clone-API/pkg/util/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const refreshTokenTest = "123123-asda1-123123"

var errSaveRefreshT = fmt.Errorf("could not save token")
var errGetRefreshT = fmt.Errorf("could not get token")

func TestRefreshTokenService(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, r *mock_repositories.MockRefreshTokenRepository,
		svc RefreshTokenService){
		"test creation of refresh token success":      testRefreshTokenCreate,
		"test creation refresh token db err fails":    testRTCreateDBErr,
		"test refresh token validation success":       testRefreshTokenValidation,
		"test refresh token validation db err fails":  testRTValidationDBErr,
		"test refresh token validation expired fails": testRTValidationExpired,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			rTRepo, svc, teardown := setupRTSvc(t)
			defer teardown()
			fn(t, rTRepo, svc)
		})
	}
}

func setupRTSvc(t *testing.T) (*mock_repositories.MockRefreshTokenRepository, RefreshTokenService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	rtRepo := mock_repositories.NewMockRefreshTokenRepository(ctrl)

	svc := NewRefreshTokenService(rtRepo)

	return rtRepo, svc, func() {
		svc = nil
		defer ctrl.Finish()
	}
}

func testRefreshTokenCreate(t *testing.T, rTRepo *mock_repositories.MockRefreshTokenRepository, svc RefreshTokenService) {
	t.Helper()

	ctx := context.Background()
	rTRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	rToken, commonError := svc.Create(ctx)
	require.NoError(t, commonError)
	require.NotEmpty(t, rToken)
}

func testRTCreateDBErr(t *testing.T, rTRepo *mock_repositories.MockRefreshTokenRepository, svc RefreshTokenService) {
	t.Helper()

	ctx := context.Background()

	testErr := errors.NewInternalServerError("Database error", errSaveRefreshT)
	rTRepo.EXPECT().Save(ctx, gomock.Any()).Return(testErr)

	rToken, commonError := svc.Create(ctx)
	require.Empty(t, rToken)
	require.Error(t, commonError)
}

func testRefreshTokenValidation(t *testing.T, rTRepo *mock_repositories.MockRefreshTokenRepository, svc RefreshTokenService) {
	t.Helper()

	refreshT := &model.RefreshToken{
		ID:        1,
		Token:     refreshTokenTest,
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	ctx := context.Background()
	rTRepo.EXPECT().FindByToken(ctx, gomock.Eq(refreshTokenTest)).Return(refreshT, nil)

	validateErr := svc.Validate(ctx, refreshTokenTest)
	require.NoError(t, validateErr)
}

func testRTValidationDBErr(t *testing.T, rTRepo *mock_repositories.MockRefreshTokenRepository, svc RefreshTokenService) {
	t.Helper()

	testErr := errors.NewInternalServerError("Database error", errGetRefreshT)

	ctx := context.Background()
	rTRepo.EXPECT().FindByToken(ctx, gomock.Eq(refreshTokenTest)).Return(nil, testErr)

	validateErr := svc.Validate(ctx, refreshTokenTest)
	require.Error(t, validateErr)
	require.Equal(t, testErr.Cause(), validateErr.Cause())
}

func testRTValidationExpired(t *testing.T, rTRepo *mock_repositories.MockRefreshTokenRepository, svc RefreshTokenService) {
	t.Helper()

	refreshT := &model.RefreshToken{
		ID:        1,
		Token:     refreshTokenTest,
		ExpiresAt: time.Now(),
	}

	ctx := context.Background()
	rTRepo.EXPECT().FindByToken(ctx, gomock.Eq(refreshTokenTest)).Return(refreshT, nil)

	validateErr := svc.Validate(ctx, refreshTokenTest)
	require.Error(t, validateErr)
}
