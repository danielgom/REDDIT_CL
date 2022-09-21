package services

import (
	"context"
	"testing"
	"time"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/services/mock_repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setupUserSvc(t *testing.T) (*mock_repositories.MockUserRepository, *mock_repositories.MockTokenRepository,
	UserService, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	tokenRepo := mock_repositories.NewMockTokenRepository(ctrl)

	svc := NewUserService(userRepo, tokenRepo)

	return userRepo, tokenRepo, svc, func() {
		svc = nil
		defer ctrl.Finish()
	}
}

func TestUserSvc_SignUp(t *testing.T) {
	t.Parallel()

	userRepo, tokenRepo, svc, teardown := setupUserSvc(t)
	defer teardown()

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

	tokenRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	userRepo.EXPECT().Save(ctx, gomock.Any()).Return(&m, nil)
	res, err := svc.SignUp(ctx, &rr)

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

func TestUserSvc_VerifyAccount(t *testing.T) {
	t.Parallel()

	userRepo, tokenRepo, svc, teardown := setupUserSvc(t)
	defer teardown()

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

	userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	tokenRepo.EXPECT().FindByToken(ctx, gomock.Any()).Return(verificationT, nil)

	verErr := svc.VerifyAccount(ctx, testVerToken)
	require.NoError(t, verErr)
}

func TestUserSvc_LoginByUsername(t *testing.T) {
	t.Parallel()

	userRepo, _, svc, teardown := setupUserSvc(t)
	defer teardown()

	loginReq := &internal.LoginRequest{
		UserOrEmail: "Daniel",
		Password:    "TestPass1@",
	}

	ctx := context.Background()

	userRepo.EXPECT().FindByUsername(ctx, gomock.Any()).Return(&model.User{
		Username: "daniel",
		Email:    "dga_355@hotmail.com",
		Password: "TestPass1@",
	}, nil)

	loginResponse, commonError := svc.Login(ctx, loginReq)

	require.Error(t, commonError)
	require.Nil(t, loginResponse)
}
