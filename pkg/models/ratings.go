package models

type Rating struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
	CafeID int32 `json:"cafe_id"`
	Rating int32 `json:"rating"`
}
