package models

type Cafe struct {
	ID           int32             `json:"id"`
	OwnerID      int32             `json:"owner_id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	OpeningTime  int8              `json:"opening_time"`
	ClosingTime  int8              `json:"closing_time"`
	Menus        []MenuItem        `json:"menus"`
	Comments     []Comment         `json:"comments"`
	Rating       float64           `json:"rating"`
	Images       []string          `json:"photos"`
	Events       []Event           `json:"events"`
	Reservations []Reservation     `json:"reservations"`
	Capacity     int32             `json:"capacity"`
	ContactInfo  ContactInfo       `json:"contact_info"`
	Categories   []CafeCategory    `json:"categories"`
	Amenities    []AmenityCategory `json:"amenities"`
}

type ContactInfo struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Province int    `json:"province"`
	City     int    `json:"city"`
	Address  string `json:"address"`
}
