package db

import (
	"context"
	"testing"
	"time"

	"RD-Clone-API/pkg/db/mock_db"
	"RD-Clone-API/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
)

func TestTokenRepository(t *testing.T) {
	t.Parallel()
	for scenario, fn := range map[string]func(t *testing.T, db *mock_db.MockDBConn, tkRepo TokenRepository){
		"test user sing up success": testSaveSuccess,
	} {
		fn := fn
		t.Run(scenario, func(t *testing.T) {
			t.Parallel()
			DBConn, repo, teardown := setupUserSvc(t)
			defer teardown()
			fn(t, DBConn, repo)
		})
	}
}

func setupUserSvc(t *testing.T) (*mock_db.MockDBConn, TokenRepository, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	DBConn := mock_db.NewMockDBConn(ctrl)

	repo := NewTokenRepository(DBConn)

	return DBConn, repo, func() {
		repo = nil
		defer ctrl.Finish()
	}
}

func testSaveSuccess(t *testing.T, conn *mock_db.MockDBConn, repo TokenRepository) {
	t.Helper()

	ctx := context.Background()
	tag := pgconn.CommandTag{}
	conn.EXPECT().Exec(ctx, gomock.Any(), gomock.Any()).Return(tag, nil)

	token := &model.VerificationToken{
		Token:      "123123",
		User:       &model.User{ID: 1},
		ExpiryDate: time.Now(),
	}

	saveErr := repo.Save(ctx, token)
	require.NoError(t, saveErr)
}
