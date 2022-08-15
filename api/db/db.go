package db

import "time"

type Mood struct {
	UserID    int       `json:"user_id"`
	Mood      int       `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
}
