package models

type MenuItem struct {
	ID          int32            `json:"id"`
	CafeID      int32            `json:"cafe_id"`
	Name        string           `json:"name"`
	Price       float64          `json:"price"`
	Category    MenuItemCategory `json:"category"`
	Ingredients []string         `json:"ingredients"`
	ImageID     string           `json:"image_id"`
}
