package deepseek

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"helixlife/internal/models"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) Chat(request models.DeepSeekRequest) (*models.DeepSeekResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败: %s", string(body))
	}

	var response models.DeepSeekResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	return &response, nil
}

// ChatStream 调用DS流式接口
func (c *Client) ChatStream(request models.DeepSeekRequest) (*http.Response, error) {
	request.Stream = true
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}
	req, err := http.NewRequest("POST", c.BaseURL+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	return c.HTTPClient.Do(req)
}

// ProcessStream 处理流式响应
func (c *Client) ProcessStream(resp *http.Response, callback func(string) error) error {
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取响应失败: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 处理SSE格式
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			if data == "[DONE]" {
				break
			}
			if err := callback(data); err != nil {
				return err
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}
