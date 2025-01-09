package models

import "time"

// User represents the user model
type User struct {
	BaseModel
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	PhoneNumber   *string   `json:"phone_number" db:"phone_number"`
	Password      string    `json:"password" db:"password"`
	EmailVerified time.Time `json:"email_verified" db:"email_verified"`
	Image         string    `json:"image" db:"image"`
	Sessions      []Session `json:"sessions"`
}

// Session represents a session for a user
type Session struct {
	BaseModel
	UserID       int       `json:"user_id" db:"user_id"`
	Expires      time.Time `json:"expires" db:"expires"`
	SessionToken string    `json:"session_token" db:"session_token"`
	User         User      `json:"user"`
}

// VerificationToken represents a token for verifying user actions
type VerificationToken struct {
	BaseModel
	ID      int64     `json:"id" db:"id"`
	Expires time.Time `json:"expires" db:"expires"`
	Token   string    `json:"token" db:"token"`
}
