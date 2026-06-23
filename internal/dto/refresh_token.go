package dto

import "time"

type RefreshToken struct {
	UserId string        `json:"user_id"`
	Token  string        `json:"token"`
	Lt     time.Duration `json:"lt"`
	Jti    string        `json:"jti"`
}
