package models

import "time"

type Cafe struct {
	ID           int32          `json:"id"`
	OwnerID      int32          `json:"owner_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	OpeningTime  *time.Time     `json:"opening_time"`
	ClosingTime  *time.Time     `json:"closing_time"`
	Menus        []MenuItem     `json:"menus"`
	Comments     []Comment      `json:"comments"`
	Rating       float64        `json:"rating"`
	Images       []string       `json:"photos"`
	Events       []Event        `json:"events"`
	Reservations []Reservation  `json:"reservations"`
	Capacity     int32          `json:"capacity"`
	ContactInfo  ContactInfo    `json:"contact_info"`
	Categories   []CafeCategory `json:"categories"`
}

type MenuItem struct {
	ID       int32            `json:"id"`
	Name     string           `json:"name"`
	Price    float64          `json:"price"`
	Category MenuItemCategory `json:"category"`
}

type ContactInfo struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Location string `json:"location"`
	Province string `json:"province"`
	City     string `json:"city"`
	Address  string `json:"address"`
}
