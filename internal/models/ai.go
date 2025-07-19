package models

// FunctionInfo 函数信息
type FunctionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
}

// TranslateRequest 翻译请求
type TranslateRequest struct {
	Text string `json:"text" binding:"required"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type TranslateResponse struct {
	OriginalText   string `json:"original_text"`
	TranslatedText string `json:"translated_text"`
	From           string `json:"from,omitempty"`
	To             string `json:"to,omitempty"`
}

// 总结请求
type SummaryRequest struct {
	Text string `json:"text" binding:"required"`
}

type SummaryResponse struct {
	OriginalText string `json:"original_text"`
	Summary      string `json:"summary"`
	Length       int    `json:"length"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
	Delta   Message `json:"delta,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type DeepSeekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type DeepSeekResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type StreamResponse struct {
	Data  string `json:"data"`
	Event string `json:"event"`
}

type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
}

type StreamRequest struct {
	Text      string `json:"text" binding:"required"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	MaxLength int    `json:"max_length,omitempty"`
}
