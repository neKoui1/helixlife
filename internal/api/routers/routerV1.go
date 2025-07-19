package routers

import (
	"helixlife/internal/api/controllers.go"
	"helixlife/internal/config"
	"helixlife/internal/services"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	cfg := config.Load()

	// 初始化服务
	aiService := services.NewAIService(cfg.DeepSeek.APIKey, cfg.DeepSeek.BaseURL)
	taskService := services.NewTaskService()
	// 初始化控制器
	aiController := controllers.NewAIController(aiService)
	taskController := controllers.NewTaskController(taskService, aiService)
	streamController := controllers.NewStreamController(aiService)
	// 初始化路由
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/functions", aiController.GetFunctions)
		apiV1.POST("/translate/zh2en", aiController.TranslateZh2En)
		apiV1.POST("/translate/en2zh", aiController.TranslateEn2Zh)
		apiV1.POST("/summary", aiController.Summary)

		// 异步任务管理
		apiV1.POST("/tasks", taskController.CreateTask)
		apiV1.GET("/tasks/:task_id", taskController.GetTask)

		// 流式处理
		apiV1.POST("/stream/translate", streamController.StreamTranslate)
		apiV1.POST("/stream/summary", streamController.StreamSummary)
	}

	return router
}
