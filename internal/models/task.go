package models

import "time"

// 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// 任务类型
type TaskType string

const (
	TaskTypeTranslate TaskType = "translate"
	TaskTypeSummary   TaskType = "summary"
)

type Task struct {
	ID        string     `json:"id"`
	Type      TaskType   `json:"type"`
	Status    TaskStatus `json:"status"`
	Request   any        `json:"request"`
	Response  any        `json:"response,omitempty"`
	Error     string     `json:"error,omitempty"`
	Progress  int        `json:"progress,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}

type CreateTaskRequest struct {
	Type    TaskType `json:"type" binding:"required"`
	Request any      `json:"request" binding:"required"`
}

type TaskResponse struct {
	ID        string     `json:"id"`
	Type      TaskType   `json:"type"`
	Status    TaskStatus `json:"status"`
	Response  any        `json:"response,omitempty"`
	Error     string     `json:"error,omitempty"`
	Progress  int        `json:"progress,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
}
