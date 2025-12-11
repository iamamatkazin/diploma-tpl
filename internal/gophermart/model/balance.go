package model

import "time"

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
