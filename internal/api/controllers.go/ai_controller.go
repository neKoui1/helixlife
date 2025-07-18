package controllers

import (
	"fmt"
	"helixlife/internal/models"
	"helixlife/internal/services"
	"helixlife/internal/utils"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiService *services.AIService
}

func NewAIController(aiService *services.AIService) *AIController {
	return &AIController{
		aiService: aiService,
	}
}

// GetFunctions 获取所有功能列表
func (h *AIController) GetFunctions(c *gin.Context) {
	functions := h.aiService.GetFunctions()
	utils.Success(c, functions)
}

// TranslateZh2En 中译英
func (h *AIController) TranslateZh2En(c *gin.Context) {
	var req models.TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, fmt.Sprintf("请求参数错误: %v", err.Error()))
		return
	}
	req.From = "zh"
	req.To = "en"
	resp, err := h.aiService.Translate(&req)
	if err != nil {
		utils.Error(c, fmt.Sprintf("翻译失败: %v", err.Error()))
		return
	}
	utils.Success(c, resp)
}

// TranslateEn2Zh 英译中
func (h *AIController) TranslateEn2Zh(c *gin.Context) {
	var req models.TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, fmt.Sprintf("请求参数错误: %v", err.Error()))
		return
	}
	req.From = "en"
	req.To = "zh"
	resp, err := h.aiService.Translate(&req)
	if err != nil {
		utils.Error(c, fmt.Sprintf("翻译失败: %v", err.Error()))
		return
	}
	utils.Success(c, resp)
}

// Summary 文本总结
func (h *AIController) Summary(c *gin.Context) {
	var req models.SummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, fmt.Sprintf("请求参数错误: %v", err.Error()))
		return
	}
	resp, err := h.aiService.Summary(&req)
	if err != nil {
		utils.Error(c, fmt.Sprintf("文本总结失败: %v", err.Error()))
		return
	}
	utils.Success(c, resp)
}
