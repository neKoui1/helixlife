package controllers

import (
	"fmt"
	"helixlife/internal/models"
	"helixlife/internal/services"
	"helixlife/internal/utils"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService *services.TaskService
	aiService   *services.AIService
}

func NewTaskController(taskService *services.TaskService, aiService *services.AIService) *TaskController {
	return &TaskController{
		taskService: taskService,
		aiService:   aiService,
	}
}

// CreateTask 创建异步任务
func (h *TaskController) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, fmt.Sprintf("请求参数错误: %v", err.Error()))
		return
	}
	var request any

	switch req.Type {
	case models.TaskTypeTranslate:
		if translateMap, ok := req.Request.(map[string]any); ok {
			request = &models.TranslateRequest{
				Text: translateMap["text"].(string),
				From: translateMap["from"].(string),
				To:   translateMap["to"].(string),
			}
		} else {
			utils.Error(c, "翻译请求参数错误")
			return
		}
	case models.TaskTypeSummary:
		if summaryMap, ok := req.Request.(map[string]any); ok {
			request = &models.SummaryRequest{
				Text: summaryMap["text"].(string),
			}
		} else {
			utils.Error(c, "总结请求参数错误")
			return
		}
	default:
		utils.Error(c, "不支持的任务类型")
		return
	}
	task := h.taskService.CreateTask(req.Type, request)

	// 异步执行任务
	go h.taskService.ExecuteTask(task.ID, h.aiService)
	utils.Success(c, gin.H{
		"id":      task.ID,
		"status":  task.Status,
		"message": "任务已创建，请使用task_id查询进度",
	})
}

func (h *TaskController) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")
	task, exists := h.taskService.GetTask(taskID)
	if !exists {
		utils.Error(c, "任务不存在")
		return
	}

	taskResp := models.TaskResponse{
		ID:        task.ID,
		Type:      task.Type,
		Status:    task.Status,
		Response:  task.Response,
		Error:     task.Error,
		Progress:  task.Progress,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		StartedAt: task.StartedAt,
		EndedAt:   task.EndedAt,
	}
	utils.Success(c, taskResp)
}
