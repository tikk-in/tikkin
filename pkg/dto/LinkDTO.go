package dto

import (
	"time"
)

// LinkDTO Link information
// @Description Link information
type LinkDTO struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Banned      bool       `json:"banned"`
	ExpireAt    *time.Time `json:"expire_at"`
	TargetUrl   string     `json:"target_url"`
	CreatedAt   *time.Time `json:"created_at" example:"2024-04-22T10:10:10Z"`
	UpdatedAt   *time.Time `json:"updated_at" example:"2024-04-22T10:10:10Z"`
	Visits      int64      `json:"visits"`
}
