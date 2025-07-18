package services

import (
	"fmt"
	"helixlife/internal/models"
	"helixlife/pkg/deepseek"
)

type AIService struct {
	Client *deepseek.Client
}

func NewAIService(apiKey, baseURL string) *AIService {
	return &AIService{
		Client: deepseek.NewClient(apiKey, baseURL),
	}
}

func (s *AIService) GetFunctions() []models.FunctionInfo {
	return []models.FunctionInfo{
		{
			ID:          "zh2en",
			Name:        "中译英",
			Description: "将中文文本翻译为英文",
			Endpoint:    "/api/v1/translate/zh2en",
			Method:      "POST",
		},
		{
			ID:          "en2zh",
			Name:        "英译中",
			Description: "将英文文本翻译为中文",
			Endpoint:    "/api/v1/translate/en2zh",
			Method:      "POST",
		},
		{
			ID:          "summary",
			Name:        "文本总结",
			Description: "对文本进行总结",
			Endpoint:    "/api/v1/summary",
			Method:      "POST",
		},
	}
}

// Translate 翻译文本
func (s *AIService) Translate(req *models.TranslateRequest) (*models.TranslateResponse, error) {
	var prompt string
	if req.From == "zh" && req.To == "en" {
		prompt = fmt.Sprintf("请将以下中文翻译为英文：\n%s", req.Text)

	} else if req.From == "en" && req.To == "zh" {
		prompt = fmt.Sprintf("Please translate the following English text into Chinese:\n%s", req.Text)
	} else {
		prompt = fmt.Sprintf("请将以下文本从%s语言翻译为%s语言：\n%s", req.From, req.To, req.Text)
	}

	dsReq := &models.DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []models.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 1.3,
		MaxTokens:   1000,
	}

	resp, err := s.Client.Chat(*dsReq)
	if err != nil {
		return nil, fmt.Errorf("翻译失败: %v", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("翻译失败: 未返回结果")
	}

	translatedText := resp.Choices[0].Message.Content
	return &models.TranslateResponse{
		OriginalText:   req.Text,
		TranslatedText: translatedText,
		From:           req.From,
		To:             req.To,
	}, nil
}

// Summary 文本总结
func (s *AIService) Summary(req *models.SummaryRequest) (*models.SummaryResponse, error) {
	prompt := fmt.Sprintf("请对以下文本进行总结：\n%s", req.Text)
	dsReq := &models.DeepSeekRequest{
		Model: "deepseek-chat",
		Messages: []models.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	resp, err := s.Client.Chat(*dsReq)
	if err != nil {
		return nil, fmt.Errorf("总结失败: %v", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("总结失败: 未返回结果")
	}

	summary := resp.Choices[0].Message.Content
	return &models.SummaryResponse{
		OriginalText: req.Text,
		Summary:      summary,
		Length:       len(summary),
	}, nil
}
