package models

type Sex int

const (
	Invalid Sex = iota
	Unknow
	Male
	Female
)

type User struct {
	ID        int32  `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     int64  `json:"phone"`
	Sex       Sex    `json:"sex"`
}
