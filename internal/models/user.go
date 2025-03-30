package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	GitHubID  string    `json:"github_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

// ToProfile converts a User to a UserProfile
func (u *User) ToProfile() *UserProfile {
	return &UserProfile{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
