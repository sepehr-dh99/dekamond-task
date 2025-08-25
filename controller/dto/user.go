package dto

import "time"

type UserResponse struct {
	Phone        string    `json:"phone" example:"09123456789"`
	RegisteredAt time.Time `json:"registered_at" example:"2025-08-25T12:00:00Z"`
}
