package models

import "time"

type Status struct {
	ID           int        `json:"-"`
	ExternalID   string     `json:"external_id"`
	Name         string     `json:"name"`
	Color        *string    `json:"color,omitempty"`
	Position     int        `json:"position"`
	ActiveStatus int        `json:"active_status"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	ModifiedAt   *time.Time `json:"modified_at,omitempty"`
	ModifiedBy   *string    `json:"modified_by,omitempty"`
}

type StatusResponse struct {
	ExternalID string     `json:"external_id"`
	Name       string     `json:"name"`
	Color      *string    `json:"color,omitempty"`
	Position   int        `json:"position"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}

type StatusRequest struct {
	Name     string  `json:"name" binding:"required"`
	Color    *string `json:"color"`
	Position *int    `json:"position"` // Optional input, default 0
}
