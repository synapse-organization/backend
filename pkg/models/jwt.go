package models

import "time"

type JWTToken struct {
	TokenID   int32     `json:"token_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	UserID    int32     `json:"user_id"`
	Role      Role      `json:"role"`
}
