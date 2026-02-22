package models

import "time"

type Workspace struct {
	ID              int        `json:"-"`
	ExternalID      string     `json:"external_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	OwnerID         int        `json:"-"`
	OwnerExternalID string     `json:"-"` // Not output as json, used for mapping
	ActiveStatus    int        `json:"active_status"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       *string    `json:"created_by,omitempty"`
	ModifiedAt      *time.Time `json:"modified_at,omitempty"`
	ModifiedBy      *string    `json:"modified_by,omitempty"`
}

type WorkspaceResponse struct {
	ExternalID      string     `json:"external_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	OwnerExternalID string     `json:"owner_external_id"`
	CreatedAt       time.Time  `json:"created_at"`
	ModifiedAt      *time.Time `json:"modified_at,omitempty"`
}

type WorkspaceRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}
