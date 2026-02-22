package models

import "time"

type Board struct {
	ID                  int        `json:"-"`
	ExternalID          string     `json:"external_id"`
	WorkspaceID         int        `json:"-"`
	WorkspaceExternalID string     `json:"-"` // Not output as json, used for mapping
	CreatedByID         *int       `json:"-"`
	Name                string     `json:"name"`
	Description         *string    `json:"description,omitempty"`
	ActiveStatus        int        `json:"active_status"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *string    `json:"created_by,omitempty"`
	ModifiedAt          *time.Time `json:"modified_at,omitempty"`
	ModifiedBy          *string    `json:"modified_by,omitempty"`
}

type BoardResponse struct {
	ExternalID          string     `json:"external_id"`
	WorkspaceExternalID string     `json:"workspace_external_id"`
	Name                string     `json:"name"`
	Description         *string    `json:"description,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	ModifiedAt          *time.Time `json:"modified_at,omitempty"`
}

type BoardRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}
