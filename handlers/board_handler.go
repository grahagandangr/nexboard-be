package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/services"
	"github.com/grahagandangr/nexboard-be/utils"
)

type BoardHandler struct {
	boardService *services.BoardService
}

func NewBoardHandler(boardService *services.BoardService) *BoardHandler {
	return &BoardHandler{boardService: boardService}
}

// CreateWorkspaceBoard handles creating a new board
func (h *BoardHandler) CreateWorkspaceBoard(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	var req models.BoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	board, err := h.boardService.CreateBoard(userExtID.(string), workspaceExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 201, board)
}

// GetWorkspaceBoards lists boards within a workspace
func (h *BoardHandler) GetWorkspaceBoards(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	workspaceExtID := c.Param("external_id")

	boards, err := h.boardService.GetWorkspaceBoards(userExtID.(string), workspaceExtID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, boards)
}

// GetBoard returns specific board details
func (h *BoardHandler) GetBoard(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	boardExtID := c.Param("external_id")

	board, err := h.boardService.GetBoard(userExtID.(string), boardExtID)
	if err != nil {
		utils.ErrorResponse(c, 404, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, board)
}

// UpdateBoard updates board details
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	boardExtID := c.Param("external_id")

	var req models.BoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	board, err := h.boardService.UpdateBoard(userExtID.(string), boardExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, board)
}

// DeleteBoard deletes a board
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	boardExtID := c.Param("external_id")

	if err := h.boardService.DeleteBoard(userExtID.(string), boardExtID); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "board deleted successfully"})
}
