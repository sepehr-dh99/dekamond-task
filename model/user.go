package model

import "time"

type User struct {
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
}
