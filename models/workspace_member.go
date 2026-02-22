package models

import "time"

type WorkspaceMember struct {
	ID          int        `json:"-"`
	WorkspaceID int        `json:"-"`
	UserID      int        `json:"-"`
	Role        string     `json:"role"`
	JoinedAt    time.Time  `json:"joined_at"`
	CreatedBy   *string    `json:"created_by,omitempty"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
	ModifiedBy  *string    `json:"modified_by,omitempty"`
}

type WorkspaceMemberResponse struct {
	UserExternalID string `json:"user_external_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

type InviteMemberRequest struct {
	UserExternalID string `json:"user_external_id" binding:"required"`
	Role           string `json:"role" binding:"required,oneof=owner admin member"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner admin member"`
}
