// Package model contains all models which are going to be inserted into the DB.
package model

import "time"

// User is the user which is going to be saved into the DB.
type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Enabled   bool
}

// VerificationToken is the starting token that we provide to the user in order to activate their account.
type VerificationToken struct {
	ID         int
	Token      string
	User       *User
	ExpiryDate time.Time
}

// RefreshToken is the token used to create a new JWT token if it has already expired.
type RefreshToken struct {
	ID        int
	Token     string
	ExpiresAt time.Time
}

// Subreddit represents a subreddit.
type Subreddit struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UserID      int
}

// Post represents a post linked to a subreddit.
type Post struct {
	ID          int
	Name        string
	URL         string
	Description string
	VoteCount   int
	UserID      int
	SubredditID int
	CreatedAt   time.Time
}
