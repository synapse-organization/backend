package models

import "time"

type Reservation struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"user_id"`
	CafeID        int32     `json:"cafe_id"`
	TransactionID string    `json:"transaction_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	People        int32     `json:"people"`
}
