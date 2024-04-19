package model

import "time"

type Link struct {
	ID          int        `json:"id"`
	UserId      int        `json:"-"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Banned      bool       `json:"banned"`
	ExpireAt    *time.Time `json:"expire_at"`
	TargetUrl   string     `json:"target_url"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
