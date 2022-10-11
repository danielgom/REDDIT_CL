// Package services contains all core logic from this API.
package services

import (
	"context"
	mErr "errors"
	"net/http"
	"net/mail"
	"time"

	"RD-Clone-API/pkg/db"
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/security"
	"RD-Clone-API/pkg/util"
	"RD-Clone-API/pkg/util/errors"
	"github.com/google/uuid"
)

const verificationTokenExpiration = 24

var errInvalidLoginRequest = mErr.New("please provide a username or an email")

type userSvc struct {
	userDB  db.UserRepository
	tokenDB db.TokenRepository
	rtSvc   RefreshTokenService
}

// NewUserService returns a new instance of user service.
func NewUserService(uR db.UserRepository, tR db.TokenRepository, rTSvc RefreshTokenService) UserService {
	return &userSvc{userDB: uR, tokenDB: tR, rtSvc: rTSvc}
}

// SignUp executes core logic in order to save the user and generate its verification token for the first time.
func (u *userSvc) SignUp(ctx context.Context, req *internal.RegisterRequest) (*internal.RegisterResponse,
	errors.CommonError) {
	pass, err := security.Hash(req.Password)
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

	user, saveErr := u.userDB.Save(ctx, user)
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

// VerifyAccount verifies the account.
func (u *userSvc) VerifyAccount(ctx context.Context, tStr string) errors.CommonError {
	token, verErr := u.tokenDB.FindByToken(ctx, tStr)
	if verErr != nil {
		return verErr
	}

	token.User.Enabled = true
	token.User.UpdatedAt = time.Now()

	updateErr := u.userDB.Update(ctx, token.User)

	if updateErr != nil {
		return updateErr
	}

	return nil
}

// Login validates username/email and password returning a JWT token and a refresh token with expiration.
func (u *userSvc) Login(ctx context.Context, loginReq *internal.LoginRequest) (*internal.LoginResponse, errors.CommonError) {
	var user *model.User
	var findErr errors.CommonError

	_, err := mail.ParseAddress(loginReq.UserOrEmail)
	if err != nil {
		user, findErr = u.userDB.FindByUsername(ctx, loginReq.UserOrEmail)
	} else {
		user, findErr = u.userDB.FindByEmail(ctx, loginReq.UserOrEmail)
	}

	if findErr != nil {
		return nil, findErr
	}

	validPass := security.CheckHash(loginReq.Password, user.Password)
	if !validPass {
		return nil, errors.NewUnauthorisedError("invalid password")
	}

	JWT, expDate, err := security.GenerateTokenWithExp(user.Username)
	if err != nil {
		return nil, errors.NewInternalServerError("internal error", err)
	}

	refreshToken, createRTErr := u.rtSvc.Create(ctx)
	if createRTErr != nil {
		return nil, createRTErr
	}

	return internal.BuildLoginResponse(user.Username, user.Email, JWT, refreshToken, expDate), nil
}

func (u *userSvc) RefreshToken(ctx context.Context, rtReq *internal.RefreshTokenRequest) (*internal.RefreshTokenResponse,
	errors.CommonError) {
	validError := u.rtSvc.Validate(ctx, rtReq.RefreshToken)
	if validError != nil {
		return nil, validError
	}

	JWT, expDate, err := security.GenerateTokenWithExp(rtReq.Username)
	if err != nil {
		return nil, errors.NewInternalServerError("internal error", err)
	}

	refreshToken, createRTErr := u.rtSvc.Create(ctx)
	if createRTErr != nil {
		return nil, createRTErr
	}

	return &internal.RefreshTokenResponse{
		Username:     rtReq.Username,
		Token:        JWT,
		RefreshToken: refreshToken,
		ExpiresAt:    expDate,
	}, nil
}

func (u *userSvc) generateVerificationToken(ctx context.Context, user *model.User) (string, errors.CommonError) {
	token := uuid.New().String()

	var vToken model.VerificationToken

	vToken.Token = token
	vToken.User = user
	vToken.ExpiryDate = time.Now().Add(time.Hour * verificationTokenExpiration)

	saveTknErr := u.tokenDB.Save(ctx, &vToken)
	if saveTknErr != nil {
		return "", saveTknErr
	}

	return token, nil
}
