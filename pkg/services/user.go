// Package services contains all core logic from this API.
package services

import (
	"context"
	"net/http"
	"time"

	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util"
	"RD-Clone-API/pkg/util/errors"
	"github.com/google/uuid"
)

const verificationTokenExpiration = 24

type userSvc struct {
	UserDB  db.UserRepository
	TokenDB db.TokenRepository
}

// NewUserService returns a new instance of user service.
func NewUserService(uR db.UserRepository, tR db.TokenRepository) UserService {
	return &userSvc{UserDB: uR, TokenDB: tR}
}

// SignUp executes core logic in order to save the user and generate its verification token for the first time.
func (u *userSvc) SignUp(ctx context.Context, req *internal.RegisterRequest) (*internal.RegisterResponse,
	errors.CommonError) {
	pass, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewRestError("password not encrypted",
			http.StatusInternalServerError, "Internal server error", err)
	}

	currentTime := time.Now().Local()

	user := new(model.User)

	user.Username = req.Username
	user.Email = req.Email
	user.Password = pass
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	user, saveErr := u.UserDB.Save(ctx, user)
	if saveErr != nil {
		return nil, saveErr
	}

	token, tknErr := u.generateVerificationToken(ctx, user)
	if tknErr != nil {
		return nil, tknErr
	}

	go util.SendRegistrationEmail(token, user.Email)

	return internal.BuildRegisterResponse(user), nil
}

func (u *userSvc) generateVerificationToken(ctx context.Context, user *model.User) (string, errors.CommonError) {
	token := uuid.New().String()

	var vToken model.VerificationToken

	vToken.Token = token
	vToken.User = user
	vToken.ExpiryDate = time.Now().Add(time.Hour * verificationTokenExpiration)

	saveTknErr := u.TokenDB.Save(ctx, &vToken)
	if saveTknErr != nil {
		return "", saveTknErr
	}

	return token, nil
}
