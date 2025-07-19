package controllers

import (
	"encoding/json"
	"helixlife/internal/models"
	"helixlife/internal/services"
	"helixlife/internal/utils"
	"io"

	"github.com/gin-gonic/gin"
)

type StreamController struct {
	aiService *services.AIService
}

func NewStreamController(aiService *services.AIService) *StreamController {
	return &StreamController{aiService: aiService}
}

// StreamTranslate 流式翻译
func (h *StreamController) StreamTranslate(c *gin.Context) {
	var req models.StreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误"+err.Error())
		return
	}
	// SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Hedaers", "Cache-Control")

	responseChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		defer close(responseChan)
		defer close(errorChan)
		err := h.aiService.TranslateStream(&req, func(chunk string) error {
			responseChan <- chunk
			return nil
		})
		if err != nil {
			errorChan <- err
		}
	}()

	// 发送流式响应
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk := <-responseChan:
			// 发送数据
			data := gin.H{
				"type": "chunk",
				"data": chunk,
			}
			jsonData, _ := json.Marshal(data)
			c.SSEvent("message", string(jsonData))
			return true
		case err := <-errorChan:
			errorData := gin.H{
				"type":  "error",
				"error": err.Error(),
			}
			jsonData, _ := json.Marshal(errorData)
			c.SSEvent("message", string(jsonData))
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func (h *StreamController) StreamSummary(c *gin.Context) {
	var req models.StreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "请求参数错误"+err.Error())
		return
	}
	// SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	responseChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		defer close(responseChan)
		defer close(errorChan)
		err := h.aiService.SummaryStream(&req, func(chunk string) error {
			responseChan <- chunk
			return nil
		})
		if err != nil {
			errorChan <- err
		}
	}()

	// 发送流式响应
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk := <-responseChan:
			// 发送数据
			data := gin.H{
				"type": "chunk",
				"data": chunk,
			}
			jsonData, _ := json.Marshal(data)
			c.SSEvent("message", string(jsonData))
			return true
		case err := <-errorChan:
			errorData := gin.H{
				"type":  "error",
				"error": err.Error(),
			}
			jsonData, _ := json.Marshal(errorData)
			c.SSEvent("message", string(jsonData))
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}
