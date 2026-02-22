package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/services"
	"github.com/grahagandangr/nexboard-be/utils"
)

type WorkspaceHandler struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceHandler(workspaceService *services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{workspaceService: workspaceService}
}

// CreateWorkspace handles creating a new workspace
func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req models.WorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	userExtID, _ := c.Get("user_external_id")

	workspace, err := h.workspaceService.CreateWorkspace(userExtID.(string), &req)
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 201, workspace)
}

// GetUserWorkspaces lists all workspaces the user is a member of
func (h *WorkspaceHandler) GetUserWorkspaces(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")

	workspaces, err := h.workspaceService.GetUserWorkspaces(userExtID.(string))
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, workspaces)
}

// GetWorkspace gets a specific workspace detail
func (h *WorkspaceHandler) GetWorkspace(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	workspace, err := h.workspaceService.GetWorkspace(userExtID.(string), workspaceExtID)
	if err != nil {
		utils.ErrorResponse(c, 404, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, workspace)
}

// UpdateWorkspace updates a workspace
func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	var req models.WorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	workspace, err := h.workspaceService.UpdateWorkspace(userExtID.(string), workspaceExtID, &req)
	if err != nil {
		// Use 403 or 400 depending on exact error parsing. 400 is fine as catch all.
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, workspace)
}

// DeleteWorkspace deletes a workspace
func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	if err := h.workspaceService.DeleteWorkspace(userExtID.(string), workspaceExtID); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "workspace deleted successfully"})
}

// GetMembers lists all members of a workspace
func (h *WorkspaceHandler) GetMembers(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	members, err := h.workspaceService.GetMembers(userExtID.(string), workspaceExtID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, members)
}

// InviteMember adds a new member to a workspace
func (h *WorkspaceHandler) InviteMember(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	var req models.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	if err := h.workspaceService.InviteMember(userExtID.(string), workspaceExtID, &req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 201, gin.H{"message": "member added successfully"})
}

// UpdateMemberRole changes a member's role
func (h *WorkspaceHandler) UpdateMemberRole(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")
	targetUserExtID := c.Param("user_ext_id")

	var req models.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	if err := h.workspaceService.UpdateMemberRole(userExtID.(string), workspaceExtID, targetUserExtID, &req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "member role updated successfully"})
}

// RemoveMember removes a member from a workspace
func (h *WorkspaceHandler) RemoveMember(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")
	targetUserExtID := c.Param("user_ext_id")

	if err := h.workspaceService.RemoveMember(userExtID.(string), workspaceExtID, targetUserExtID); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "member removed successfully"})
}
