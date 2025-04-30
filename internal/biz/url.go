package biz

import "time"

// Url is a URL entity.
type Url struct {
	ID        int64
	LongURL   string
	ShortCode string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time // nil means never expires
}
