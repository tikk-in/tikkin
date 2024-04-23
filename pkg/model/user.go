package model

import "time"

type User struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Verified          bool      `json:"verified"`
	VerificationToken *string   `json:"verification_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
