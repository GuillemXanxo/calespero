package domain

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"-"` // Password is not exposed in JSON
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"created_at"`
	LastConnection time.Time `json:"last_connection"`
	Orders         []Order   `json:"orders"`
}
