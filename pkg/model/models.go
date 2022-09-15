package model

import "time"

type User struct {
	Id        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

type VerificationToken struct {
	Id         int
	Token      string
	User       *User
	ExpiryDate time.Time
}
