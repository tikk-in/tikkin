package model

import "time"

type Visits struct {
	ID          string     `json:"id"`
	LinkID      int64      `json:"link_id"`
	CreatedAt   *time.Time `json:"created_at"`
	UserAgent   *string    `json:"user_agent"`
	Referrer    *string    `json:"referrer"`
	CountryCode *string    `json:"country_code"`
}
