package services

import (
	"errors"

	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/repositories"
	"github.com/grahagandangr/nexboard-be/utils"
)

type WorkspaceService struct {
	workspaceRepo *repositories.WorkspaceRepository
	userRepo      *repositories.UserRepository
}

func NewWorkspaceService(workspaceRepo *repositories.WorkspaceRepository, userRepo *repositories.UserRepository) *WorkspaceService {
	return &WorkspaceService{workspaceRepo: workspaceRepo, userRepo: userRepo}
}

// CreateWorkspace creates a new workspace and sets the user as the owner
func (s *WorkspaceService) CreateWorkspace(creatorExternalID string, req *models.WorkspaceRequest) (*models.WorkspaceResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(creatorExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	workspace := &models.Workspace{
		ExternalID:  utils.GenerateUUID(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     user.ID,
	}

	if err := s.workspaceRepo.CreateWorkspace(workspace); err != nil {
		return nil, err
	}

	return &models.WorkspaceResponse{
		ExternalID:      workspace.ExternalID,
		Name:            workspace.Name,
		Description:     workspace.Description,
		OwnerExternalID: user.ExternalID,
		CreatedAt:       workspace.CreatedAt,
		ModifiedAt:      workspace.ModifiedAt,
	}, nil
}

// GetUserWorkspaces lists all workspaces the user is a member of
func (s *WorkspaceService) GetUserWorkspaces(userExternalID string) ([]*models.WorkspaceResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	workspaces, err := s.workspaceRepo.GetWorkspacesByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var response []*models.WorkspaceResponse
	for _, w := range workspaces {
		response = append(response, &models.WorkspaceResponse{
			ExternalID:      w.ExternalID,
			Name:            w.Name,
			Description:     w.Description,
			OwnerExternalID: w.OwnerExternalID,
			CreatedAt:       w.CreatedAt,
			ModifiedAt:      w.ModifiedAt,
		})
	}
	return response, nil
}

// GetWorkspace gets a specific workspace
func (s *WorkspaceService) GetWorkspace(userExternalID, workspaceExternalID string) (*models.WorkspaceResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	// Verify user is a member
	_, err = s.workspaceRepo.GetMemberRole(w.ID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of this workspace")
	}

	return &models.WorkspaceResponse{
		ExternalID:      w.ExternalID,
		Name:            w.Name,
		Description:     w.Description,
		OwnerExternalID: w.OwnerExternalID,
		CreatedAt:       w.CreatedAt,
		ModifiedAt:      w.ModifiedAt,
	}, nil
}

// UpdateWorkspace updates a workspace (owner only)
func (s *WorkspaceService) UpdateWorkspace(userExternalID, workspaceExternalID string, req *models.WorkspaceRequest) (*models.WorkspaceResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	// Verify user is owner
	if w.OwnerID != user.ID {
		return nil, errors.New("unauthorized: only owner can update workspace")
	}

	w.Name = req.Name
	w.Description = req.Description

	if err := s.workspaceRepo.UpdateWorkspace(w); err != nil {
		return nil, err
	}

	return &models.WorkspaceResponse{
		ExternalID:      w.ExternalID,
		Name:            w.Name,
		Description:     w.Description,
		OwnerExternalID: w.OwnerExternalID,
		CreatedAt:       w.CreatedAt,
		ModifiedAt:      w.ModifiedAt,
	}, nil
}

// DeleteWorkspace deletes a workspace (owner only)
func (s *WorkspaceService) DeleteWorkspace(userExternalID, workspaceExternalID string) error {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Verify user is owner
	if w.OwnerID != user.ID {
		return errors.New("unauthorized: only owner can delete workspace")
	}

	return s.workspaceRepo.DeleteWorkspace(w.ID)
}

// --------- Member management -----------

// GetMembers gets all members of a workspace
func (s *WorkspaceService) GetMembers(userExternalID, workspaceExternalID string) ([]*models.WorkspaceMemberResponse, error) {
	user, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	// Verify user is a member
	_, err = s.workspaceRepo.GetMemberRole(w.ID, user.ID)
	if err != nil {
		return nil, errors.New("unauthorized: not a member of this workspace")
	}

	return s.workspaceRepo.GetMembers(w.ID)
}

// InviteMember adds a member to the workspace (owner/admin only)
func (s *WorkspaceService) InviteMember(userExternalID, workspaceExternalID string, req *models.InviteMemberRequest) error {
	currentUser, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Check permission (owner or admin)
	role, err := s.workspaceRepo.GetMemberRole(w.ID, currentUser.ID)
	if err != nil {
		return errors.New("unauthorized: not a member")
	}
	if role != "owner" && role != "admin" {
		return errors.New("unauthorized: only owner or admin can invite members")
	}

	// Get target user
	targetUser, err := s.userRepo.GetUserByExternalID(req.UserExternalID)
	if err != nil {
		return errors.New("target user not found")
	}

	wm := &models.WorkspaceMember{
		WorkspaceID: w.ID,
		UserID:      targetUser.ID,
		Role:        req.Role,
	}

	err = s.workspaceRepo.AddMember(wm)
	if err != nil {
		return errors.New("could not add member (perhaps already a member?)")
	}

	return nil
}

// UpdateMemberRole updates a member's role (owner only)
func (s *WorkspaceService) UpdateMemberRole(userExternalID, workspaceExternalID, targetUserExternalID string, req *models.UpdateMemberRoleRequest) error {
	currentUser, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Verify current user is owner
	if w.OwnerID != currentUser.ID {
		return errors.New("unauthorized: only owner can update member roles")
	}

	// Get target user
	targetUser, err := s.userRepo.GetUserByExternalID(targetUserExternalID)
	if err != nil {
		return errors.New("target user not found")
	}

	// Owner cannot change their own role here
	if targetUser.ID == w.OwnerID {
		return errors.New("cannot change owner role")
	}

	return s.workspaceRepo.UpdateMemberRole(w.ID, targetUser.ID, req.Role)
}

// RemoveMember removes a member (owner/admin only)
func (s *WorkspaceService) RemoveMember(userExternalID, workspaceExternalID, targetUserExternalID string) error {
	currentUser, err := s.userRepo.GetUserByExternalID(userExternalID)
	if err != nil {
		return errors.New("user not found")
	}

	w, err := s.workspaceRepo.GetWorkspaceByExternalID(workspaceExternalID)
	if err != nil {
		return errors.New("workspace not found")
	}

	// Get current user role
	role, err := s.workspaceRepo.GetMemberRole(w.ID, currentUser.ID)
	if err != nil || (role != "owner" && role != "admin") {
		return errors.New("unauthorized: only owner or admin can remove members")
	}

	// Get target user
	targetUser, err := s.userRepo.GetUserByExternalID(targetUserExternalID)
	if err != nil {
		return errors.New("target user not found")
	}

	// Owner cannot be removed
	if targetUser.ID == w.OwnerID {
		return errors.New("cannot remove workspace owner")
	}

	return s.workspaceRepo.RemoveMember(w.ID, targetUser.ID)
}
