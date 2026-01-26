package handlers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qinyang/taskmanager/database"
	"github.com/qinyang/taskmanager/models"
	"github.com/qinyang/taskmanager/utils"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var query models.TaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Default pagination values
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	db := database.GetDB().Model(&models.Task{}).Where("user_id = ?", userID)

	// Apply filters
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Priority != "" {
		db = db.Where("priority = ?", query.Priority)
	}
	if query.Search != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	// Get total count
	var total int64
	db.Count(&total)

	// Get tasks with pagination
	var tasks []models.Task
	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&tasks).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	response := models.TaskListResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, http.StatusOK, "Tasks retrieved successfully", response)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	taskID := c.Param("id")

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	if err := database.GetDB().Where("id = ? AND user_id = ?", taskUUID, userID).First(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task retrieved successfully", task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	task := models.Task{
		UserID:      userID.(uuid.UUID),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Status:      models.StatusPending,
	}

	if err := database.GetDB().Create(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Task created successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	taskID := c.Param("id")

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	if err := database.GetDB().Where("id = ? AND user_id = ?", taskUUID, userID).First(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		updates["due_date"] = req.DueDate
	}

	if err := database.GetDB().Model(&task).Updates(updates).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	// Reload task to get updated values
	if err := database.GetDB().First(&task, taskUUID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	taskID := c.Param("id")

	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	if err := database.GetDB().Where("id = ? AND user_id = ?", taskUUID, userID).First(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}

	if err := database.GetDB().Delete(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}
