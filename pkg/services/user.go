// Package services contains all core logic from this API.
package services

import (
	"RD-Clone-API/pkg/db"
	"context"
	"time"

	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util"

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
func (u *userSvc) SignUp(ctx context.Context, req *internal.RegisterRequest) error {
	var user model.User

	pass, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	currentTime := time.Now().Local()

	user.Username = req.Username
	user.Email = req.Email
	user.Password = pass
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	_, err = u.UserDB.Save(ctx, &user)
	if err != nil {
		return err
	}

	token, err := u.generateVerificationToken(ctx, &user)
	if err != nil {
		return err
	}

	go util.SendMail("Activate Spring reddit CL account", "Thank you for signing up to Spring reddit service,"+
		" please click the link below to activate your account"+"http://localhost:8080/api/auth/accountVerification/"+token, user.Email)

	return nil
}

func (u *userSvc) generateVerificationToken(ctx context.Context, user *model.User) (string, error) {
	token := uuid.New().String()

	var vToken model.VerificationToken

	vToken.Token = token
	vToken.User = user
	vToken.ExpiryDate = time.Now().Add(time.Hour * verificationTokenExpiration)

	err := u.TokenDB.Save(ctx, &vToken)
	if err != nil {
		return "", err
	}

	return token, nil
}
