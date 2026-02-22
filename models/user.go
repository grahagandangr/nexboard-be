package models

import "time"

type User struct {
	ID           int        `json:"id"`
	ExternalID   string     `json:"external_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Password     string     `json:"password,omitempty"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	ActiveStatus int        `json:"active_status"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	ModifiedAt   *time.Time `json:"modified_at,omitempty"`
	ModifiedBy   *string    `json:"modified_by,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type UpdateProfileRequest struct {
	Name      string  `json:"name" binding:"required"`
	AvatarURL *string `json:"avatar_url"`
}
