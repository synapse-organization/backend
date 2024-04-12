package models

import "time"

type Comment struct {
	ID      int32     `json:"id"`
	UserID  int32     `json:"user_id"`
	CafeID  int32     `json:"cafe_id"`
	Comment string    `json:"text"`
	Date    time.Time `json:"date"`
}
