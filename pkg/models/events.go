package models

import "time"

type Event struct {
	ID          int32     `json:"id"`
	CafeID      int32     `json:"cafe_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	ImageID     string    `json:"image_id"`
}

type EventReservation struct {
	ID      int32 `json:"id"`
	UserID  int32 `json:"user_id"`
	EventID int32 `json:"event_id"`
}
