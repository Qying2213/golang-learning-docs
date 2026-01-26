package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"

	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	if t.Status == "" {
		t.Status = StatusPending
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	return nil
}

// DTOs
type CreateTaskRequest struct {
	Title       string       `json:"title" binding:"required,min=1,max=255"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title"`
	Description *string       `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     *time.Time    `json:"due_date"`
}

type TaskListResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

type TaskQuery struct {
	Status   TaskStatus   `form:"status"`
	Priority TaskPriority `form:"priority"`
	Page     int          `form:"page,default=1"`
	PageSize int          `form:"page_size,default=10"`
	Search   string       `form:"search"`
}
