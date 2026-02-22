package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/services"
	"github.com/grahagandangr/nexboard-be/utils"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateBoardTask processes inbound task requests
func (h *TaskHandler) CreateBoardTask(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	boardExtID := c.Param("external_id")

	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	task, err := h.taskService.CreateTask(userExtID.(string), boardExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 201, task)
}

// GetBoardTasks fetches all tasks linked to a board
func (h *TaskHandler) GetBoardTasks(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	boardExtID := c.Param("external_id")

	tasks, err := h.taskService.GetBoardTasks(userExtID.(string), boardExtID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, tasks)
}

// UpdateTask modifies a task via PUT
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	taskExtID := c.Param("external_id")

	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(userExtID.(string), taskExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, task)
}

// DeleteTask cleans a task out
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	taskExtID := c.Param("external_id")

	err := h.taskService.DeleteTask(userExtID.(string), taskExtID)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "task deleted successfully"})
}

// MoveTask handles isolated PATCH of task status
func (h *TaskHandler) MoveTask(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	taskExtID := c.Param("external_id")

	var req models.MoveTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	task, err := h.taskService.MoveTaskStatus(userExtID.(string), taskExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, task)
}

// AssignTask handles isolated PATCH of assigned user
func (h *TaskHandler) AssignTask(c *gin.Context) {
	userExtID, _ := c.Get("user_external_id")
	taskExtID := c.Param("external_id")

	var req models.AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request body")
		return
	}

	task, err := h.taskService.AssignTask(userExtID.(string), taskExtID, &req)
	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, task)
}
