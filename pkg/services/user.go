package services

import (
	"RD-Clone-API/pkg/internal"
	"RD-Clone-API/pkg/model"
	"RD-Clone-API/pkg/util"
	"context"
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	FindByUsername(context.Context, string) (*model.User, error)
	Save(context.Context, *model.User) (*model.User, error)
}

type TokenRepository interface {
	Save(context.Context, *model.VerificationToken) error
}

type UserSvc struct {
	UserDB  UserRepository
	TokenDB TokenRepository
}

func NewUserService(uR UserRepository, tR TokenRepository) *UserSvc {
	return &UserSvc{UserDB: uR, TokenDB: tR}
}

func (u *UserSvc) SignUp(ctx context.Context, req *internal.RegisterRequest) error {
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

func (u *UserSvc) generateVerificationToken(ctx context.Context, user *model.User) (string, error) {
	token := uuid.New().String()

	var vToken model.VerificationToken

	vToken.Token = token
	vToken.User = user
	vToken.ExpiryDate = time.Now().Add(time.Hour * 24)

	err := u.TokenDB.Save(ctx, &vToken)
	if err != nil {
		return "", err
	}

	return token, nil
}
