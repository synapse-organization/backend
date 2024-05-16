package models

import "time"

type Event struct {
	ID               int32     `json:"id"`
	CafeID           int32     `json:"cafe_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ImageID          string    `json:"image_id"`
	Price            float64   `json:"price"`
	Capacity         int32     `json:"capacity"`
	CurrentAttendees int32     `json:"current_attendees"`
	Reservable       bool      `json:"reservable"`
}

type EventReservation struct {
	ID      int32 `json:"id"`
	UserID  int32 `json:"user_id"`
	EventID int32 `json:"event_id"`
	TransactionID string `json:"transaction_id"`
}
