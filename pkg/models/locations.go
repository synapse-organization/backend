package models

type Location struct {
	CafeID int32   `json:"cafe_id"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}
