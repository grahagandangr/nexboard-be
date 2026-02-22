package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/services"
	"github.com/grahagandangr/nexboard-be/utils"
)

type StatusHandler struct {
	statusService *services.StatusService
}

func NewStatusHandler(statusService *services.StatusService) *StatusHandler {
	return &StatusHandler{statusService: statusService}
}

// CreateStatus handles new status
func (h *StatusHandler) CreateStatus(c *gin.Context) {
	var req models.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	status, err := h.statusService.CreateStatus(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 201, status)
}

// GetAllStatuses returns a list of all statuses
func (h *StatusHandler) GetAllStatuses(c *gin.Context) {
	statuses, err := h.statusService.GetAllStatuses()
	if err != nil {
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, statuses)
}

// GetStatus returns a single status detail
func (h *StatusHandler) GetStatus(c *gin.Context) {
	extID := c.Param("external_id")

	status, err := h.statusService.GetStatus(extID)
	if err != nil {
		utils.ErrorResponse(c, 404, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, status)
}

// UpdateStatus changes an existing status
func (h *StatusHandler) UpdateStatus(c *gin.Context) {
	extID := c.Param("external_id")

	var req models.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	status, err := h.statusService.UpdateStatus(extID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, status)
}

// DeleteStatus removes status safely
func (h *StatusHandler) DeleteStatus(c *gin.Context) {
	extID := c.Param("external_id")

	err := h.statusService.DeleteStatus(extID)
	if err != nil {
		// Translate conflict to HTTP 409
		if strings.HasPrefix(err.Error(), "conflict") {
			utils.ErrorResponse(c, 409, err.Error())
			return
		}
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "status deleted successfully"})
}
