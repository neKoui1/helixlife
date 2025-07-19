package services

import (
	"fmt"
	"helixlife/internal/models"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	tasks map[string]*models.Task
	mutex sync.RWMutex
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]*models.Task),
	}
}

func (s *TaskService) CreateTask(taskType models.TaskType, request any) *models.Task {
	now := time.Now()
	task := &models.Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Status:    models.TaskStatusPending,
		Request:   request,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.mutex.Lock()
	s.tasks[task.ID] = task
	s.mutex.Unlock()
	log.Printf("创建任务: %s, 类型: %s", task.ID, task.Type)
	return task
}

func (s *TaskService) GetTask(taskID string) (*models.Task, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	task, exists := s.tasks[taskID]
	return task, exists
}

func (s *TaskService) UpdateTask(taskID string, status models.TaskStatus, response any, err string, progress int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, exists := s.tasks[taskID]; exists {
		task.Status = status
		task.Response = response
		task.Error = err
		task.Progress = progress
		task.UpdatedAt = time.Now()

		if status == models.TaskStatusRunning && task.StartedAt == nil {
			now := time.Now()
			task.StartedAt = &now
		}

		if status == models.TaskStatusCompleted || status == models.TaskStatusFailed {
			now := time.Now()
			task.EndedAt = &now
		}
		log.Printf("更新任务: %s, 状态: %s, 进度: %d, 错误: %s", taskID, status, progress, err)
	}

}

func (s *TaskService) ExecuteTask(taskID string, aiService *AIService) {
	log.Printf("开始执行任务: %s", taskID)
	task, exists := s.GetTask(taskID)
	if !exists {
		log.Printf("任务不存在: %s", taskID)
		return
	}

	// 更新状态为运行中
	s.UpdateTask(taskID, models.TaskStatusRunning, nil, "", 10)
	// 添加延迟模拟处理时间
	time.Sleep(1 * time.Second)

	switch task.Type {
	case models.TaskTypeTranslate:
		log.Printf("执行翻译任务: %s", taskID)
		if req, ok := task.Request.(*models.TranslateRequest); ok {
			log.Printf("翻译请求: %+v", req)
			resp, err := aiService.Translate(req)
			if err != nil {
				log.Printf("翻译失败: %s", err.Error())
				s.UpdateTask(taskID, models.TaskStatusFailed, nil, err.Error(), 100)
			} else {
				log.Printf("翻译成功: %+v", resp)
				s.UpdateTask(taskID, models.TaskStatusCompleted, resp, "", 100)
			}
		} else {
			log.Printf("翻译任务类型断言失败")
			s.UpdateTask(taskID, models.TaskStatusFailed, nil, "请求数据类型错误", 100)
		}
	case models.TaskTypeSummary:
		log.Printf("执行总结任务: %s", taskID)
		if req, ok := task.Request.(*models.SummaryRequest); ok {
			log.Printf("总结请求: %+v", req)
			resp, err := aiService.Summary(req)
			if err != nil {
				log.Printf("总结失败: %s", err.Error())
				s.UpdateTask(taskID, models.TaskStatusFailed, nil, err.Error(), 100)
			} else {
				log.Printf("总结成功: %+v", resp)
				s.UpdateTask(taskID, models.TaskStatusCompleted, resp, "", 100)
			}
		} else {
			log.Printf("总结任务类型断言失败")
			s.UpdateTask(taskID, models.TaskStatusFailed, nil, "请求数据类型错误", 100)
		}
	default:
		log.Printf("不支持的任务: %s, 类型: %s", taskID, task.Type)
		s.UpdateTask(taskID, models.TaskStatusFailed, nil, fmt.Sprintf("不支持的任务类型: %s", task.Type), 100)
	}
}

func (s *TaskService) GetAllTasks() []*models.Task {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}
