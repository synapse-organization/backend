package models

import "time"

type Cafe struct {
	ID           int32         `json:"id"`
	OwnerID      int32         `json:"owner_id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	OpeningTime  time.Time     `json:"opening_time"`
	ClosingTime  time.Time     `json:"closing_time"`
	Menus        []MenuItem    `json:"menus"`
	Comments     []Comment     `json:"comments"`
	Rating       float64       `json:"rating"`
	Photos       []string      `json:"photos"`
	Events       []Event       `json:"events"`
	Reservations []Reservation `json:"reservations"`
	Capacity     int32         `json:"capacity"`
	ContactInfo  ContactInfo   `json:"contact_info"`
	Categories   []Category    `json:"categories"`
}

type MenuItem struct {
	ID       int32    `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type ContactInfo struct {
	Phone    string        `json:"phone"`
	Email    string        `json:"email"`
	Address  string        `json:"address"`
	Location time.Location `json:"location"`
}

type CategoryType string

const (
	CafeCategory CategoryType = "cafe"
	MenuCategory CategoryType = "menu"
)

type Category struct {
	ID   int32        `json:"id"`
	Name string       `json:"name"`
	Type CategoryType `json:"type"`
}
