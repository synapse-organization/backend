package models

import "time"

type Event struct {
	ID          int32     `json:"id"`
	CafeID      int32     `json:"cafe_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_date"`
	EndTime     time.Time `json:"end_date"`
	Photos      []string  `json:"photos"`
}

type EventReservation struct {
	ID      int32 `json:"id"`
	UserID  int32 `json:"user_id"`
	EventID int32 `json:"event_id"`
}
