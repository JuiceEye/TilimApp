package model

import "time"

type Subscription struct {
	ID        int64
	UserID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
