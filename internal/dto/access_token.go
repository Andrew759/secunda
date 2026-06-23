package dto

import "time"

type AccessToken struct {
	UserId string        `json:"user_id"`
	Token  string        `json:"token"`
	Lt     time.Duration `json:"lt"`
}
