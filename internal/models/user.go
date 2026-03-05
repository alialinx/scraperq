package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	APIKey       string    `json:"api_key"`
	IsActive     bool      `json:"is_active"`
	DailyLimit   int       `json:"daily_limit"`
	MonthlyLimit int       `json:"monthly_limit"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	APIKey string `json:"api_key,omitempty"`
}
