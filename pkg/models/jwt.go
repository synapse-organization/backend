package models

import "time"

type JWTToken struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       int32     `json:"user_id"`
}
