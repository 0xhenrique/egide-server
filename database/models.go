package database

import (
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"-"` // DO NOT RETURN IN THE JSON!!!
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Website struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Domain         string    `json:"domain"`
	Description    string    `json:"description,omitempty"`
	ProtectionMode string    `json:"protection_mode"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
