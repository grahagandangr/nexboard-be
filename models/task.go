package models

import "time"

type Task struct {
	ID           int        `json:"-"`
	ExternalID   string     `json:"external_id"`
	BoardID      int        `json:"-"`
	StatusID     int        `json:"-"`
	AssignedTo   *int       `json:"-"`
	CreatedByID  int        `json:"-"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	Priority     string     `json:"priority"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Position     int        `json:"position"`
	ActiveStatus int        `json:"active_status"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	ModifiedAt   *time.Time `json:"modified_at,omitempty"`
	ModifiedBy   *string    `json:"modified_by,omitempty"`
}

type TaskResponse struct {
	ExternalID        string            `json:"external_id"`
	BoardExternalID   string            `json:"board_external_id"`
	Status            TaskStatusInfo    `json:"status"`
	AssignedTo        *TaskAssigneeInfo `json:"assigned_to"`
	Title             string            `json:"title"`
	Description       *string           `json:"description,omitempty"`
	Priority          string            `json:"priority"`
	DueDate           *time.Time        `json:"due_date,omitempty"`
	Position          int               `json:"position"`
	CreatedAt         time.Time         `json:"created_at"`
	ModifiedAt        *time.Time        `json:"modified_at,omitempty"`
}

type TaskStatusInfo struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Color      *string `json:"color,omitempty"`
}

type TaskAssigneeInfo struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type TaskRequest struct {
	Title                 string     `json:"title" binding:"required"`
	Description           *string    `json:"description"`
	Priority              string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate               *time.Time `json:"due_date"`
	StatusExternalID      string     `json:"status_external_id" binding:"required"`
	AssignedToExternalID  *string    `json:"assigned_to_external_id"`
}

type MoveTaskStatusRequest struct {
	StatusExternalID string `json:"status_external_id" binding:"required"`
}

type AssignTaskRequest struct {
	AssignedToExternalID *string `json:"assigned_to_external_id"` // can be nil to unassign
}
