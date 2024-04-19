package dto

import (
	"time"
)

type LinkDTO struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Banned      bool       `json:"banned"`
	ExpireAt    *time.Time `json:"expire_at"`
	TargetUrl   string     `json:"target_url"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Visits      int64      `json:"visits"`
}
