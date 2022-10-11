package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/services/mock_repositories"
	"RD-Clone-API/pkg/services/mock_services"
	"RD-Clone-API/pkg/util/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errSaveTkn = fmt.Errorf("could not save token to db")
var errSaveUser = fmt.Errorf("could not save user to db")

func TestUserService(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, uR *mock_repositories.MockUserRepository,
		tknRepo *mock_repositories.MockTokenRepository, rTSvc *mock_services.MockRefreshTokenService, userSvc UserService){
		"test user sing up success":        testUserSignUp,
		"test verify account success":      testVerifyAccount,
		"test login with username success": testLoginByUsername,
		"test sign up db err fails":        testUserSignUpDBErr,
		"test sign up tkn db err fails":    testUserServiceSignUpTknDBErr,
		"test refresh token success":       testRefreshToken,
		"test refresh token expired fails": testRefreshTokenExpired,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			userRepo, tokenRepo, rTSvc, svc, teardown := setupUserSvc(t)
			defer teardown()
			fn(t, userRepo, tokenRepo, rTSvc, svc)
		})
	}
}

func setupUserSvc(t *testing.T) (*mock_repositories.MockUserRepository, *mock_repositories.MockTokenRepository,
	*mock_services.MockRefreshTokenService, UserService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	tokenRepo := mock_repositories.NewMockTokenRepository(ctrl)
	rTService := mock_services.NewMockRefreshTokenService(ctrl)

	svc := NewUserService(userRepo, tokenRepo, rTService)

	return userRepo, tokenRepo, rTService, svc, func() {
		svc = nil
		defer ctrl.Finish()
	}
}

func testUserSignUp(t *testing.T, uR *mock_repositories.MockUserRepository,
	tknRepo *mock_repositories.MockTokenRepository, _ *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	rr := internal.RegisterRequest{
		Username: "Daniel",
		Password: "P@ssw0rd123123",
		Email:    "dga_355@hotmail.com",
	}

	m := model.User{
		Username:  rr.Username,
		Password:  rr.Password,
		Email:     rr.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Enabled:   false,
	}

	ctx := context.Background()

	uR.EXPECT().Save(ctx, gomock.Any()).Return(&m, nil)
	tknRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	res, err := userSvc.SignUp(ctx, &rr)

	want := &internal.RegisterResponse{
		Username: rr.Username,
		Email:    rr.Email,
		Enabled:  false,
	}

	require.NoError(t, err)
	require.Equal(t, want.Username, res.Username)
	require.Equal(t, want.Enabled, res.Enabled)
	require.Equal(t, want.Email, res.Email)
}

func testUserSignUpDBErr(t *testing.T, uR *mock_repositories.MockUserRepository,
	_ *mock_repositories.MockTokenRepository, _ *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	rr := internal.RegisterRequest{
		Username: "Daniel",
		Password: "P@ssw0rd123123",
		Email:    "dga_355@hotmail.com",
	}

	ctx := context.Background()
	saveErr := errors.NewInternalServerError("could not create user", errSaveUser)

	uR.EXPECT().Save(ctx, gomock.Any()).Return(nil, saveErr)
	res, err := userSvc.SignUp(ctx, &rr)

	require.Error(t, err)
	require.Nil(t, res)
}

func testUserServiceSignUpTknDBErr(t *testing.T, uR *mock_repositories.MockUserRepository,
	tknRepo *mock_repositories.MockTokenRepository, _ *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	rr := internal.RegisterRequest{
		Username: "Daniel",
		Password: "P@ssw0rd123123",
		Email:    "dga_355@hotmail.com",
	}

	m := model.User{
		Username:  rr.Username,
		Password:  rr.Password,
		Email:     rr.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Enabled:   false,
	}

	ctx := context.Background()
	saveErr := errors.NewInternalServerError("could not create token", errSaveTkn)

	uR.EXPECT().Save(ctx, gomock.Any()).Return(&m, nil)
	tknRepo.EXPECT().Save(ctx, gomock.Any()).Return(saveErr)
	res, err := userSvc.SignUp(ctx, &rr)

	require.Error(t, err)
	require.Nil(t, res)
}

func testVerifyAccount(t *testing.T, uR *mock_repositories.MockUserRepository,
	tknRepo *mock_repositories.MockTokenRepository, _ *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	testVerToken := "abc123"

	u := model.User{
		Username:  "daniel",
		Email:     "dga_355@hotmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Enabled:   false,
	}

	verificationT := &model.VerificationToken{
		ID:    1,
		Token: testVerToken,
		User:  &u,
	}
	ctx := context.Background()

	uR.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	tknRepo.EXPECT().FindByToken(ctx, gomock.Any()).Return(verificationT, nil)

	verErr := userSvc.VerifyAccount(ctx, testVerToken)
	require.NoError(t, verErr)
}

func testLoginByUsername(t *testing.T, uR *mock_repositories.MockUserRepository,
	_ *mock_repositories.MockTokenRepository, _ *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	loginReq := &internal.LoginRequest{
		UserOrEmail: "Daniel",
		Password:    "TestPass1@",
	}

	ctx := context.Background()

	uR.EXPECT().FindByUsername(ctx, gomock.Any()).Return(&model.User{
		Username: "daniel",
		Email:    "dga_355@hotmail.com",
		Password: "TestPass1@",
	}, nil)

	loginResponse, commonError := userSvc.Login(ctx, loginReq)

	require.Error(t, commonError)
	require.Nil(t, loginResponse)
}

func testRefreshToken(t *testing.T, _ *mock_repositories.MockUserRepository,
	_ *mock_repositories.MockTokenRepository, rtSvc *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	refReq := &internal.RefreshTokenRequest{
		RefreshToken: "asd12345",
		Username:     "testuser",
	}

	want := "rt1234asd"

	ctx := context.Background()

	rtSvc.EXPECT().Validate(ctx, refReq.RefreshToken).Return(nil)
	rtSvc.EXPECT().Create(ctx).Return(want, nil)

	resp, err := userSvc.RefreshToken(ctx, refReq)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, want, resp.RefreshToken)
	require.Equal(t, refReq.Username, resp.Username)
	require.Greater(t, resp.ExpiresAt.Unix(), time.Now().Unix())

}

func testRefreshTokenExpired(t *testing.T, _ *mock_repositories.MockUserRepository,
	_ *mock_repositories.MockTokenRepository, rtSvc *mock_services.MockRefreshTokenService, userSvc UserService) {
	t.Helper()

	refReq := &internal.RefreshTokenRequest{
		RefreshToken: "asd12345",
		Username:     "testuser",
	}

	ctx := context.Background()

	rtSvc.EXPECT().Validate(ctx, refReq.RefreshToken).Return(errors.NewUnauthorisedError("refresh token expired"))

	resp, err := userSvc.RefreshToken(ctx, refReq)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, err.Message(), "refresh token expired")
}
