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

func TestUserSvc_SignUp(t *testing.T) {
	t.Parallel()
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userRepo := mock_repositories.NewMockUserRepository(ctrl)
	tokenRepo := mock_repositories.NewMockTokenRepository(ctrl)

	tokenRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	userRepo.EXPECT().Save(ctx, gomock.Any()).Return(&m, nil)

	service := NewUserService(userRepo, tokenRepo)
	err := service.SignUp(ctx, &rr)
	require.NoError(t, err)
}
