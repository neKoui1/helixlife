# helixlife
written test for helixlife

项目Github：[https://github.com/neKoui1/helixlife?tab=readme-ov-file#helixlife](https://github.com/neKoui1/helixlife?tab=readme-ov-file#helixlife)

在项目根目录下的配置文件`config.yaml`，读取使用的框架是`viper`，后端框架`gin`

```yaml
server:
  port: "8080"

deepseek:
  api_key: "your-ds-api-key"
  base_url: "https://api.deepseek.com"
```

# 实现小型AI翻译应用后端接口

## 1. 查看所有功能列表API

核心业务代码：

```go
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
```

其中，`FunctionInfo`的数据结构为：

```go
// FunctionInfo 函数信息
type FunctionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
}
```

测试API：

![image-20250718211248523](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250718211248616.png)

返回结果：

```json
{
    "code": 200,
    "message": "success",
    "data": [
        {
            "id": "zh2en",
            "name": "中译英",
            "description": "将中文文本翻译为英文",
            "endpoint": "/api/v1/translate/zh2en",
            "method": "POST"
        },
        {
            "id": "en2zh",
            "name": "英译中",
            "description": "将英文文本翻译为中文",
            "endpoint": "/api/v1/translate/en2zh",
            "method": "POST"
        },
        {
            "id": "summary",
            "name": "文本总结",
            "description": "对文本进行总结",
            "endpoint": "/api/v1/summary",
            "method": "POST"
        }
    ]
}
```

## 2. 核心翻译业务API

核心业务代码：

```go
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
```

其中，`TranslateRequest`和`TranslateResponse`的数据结构为：

```go
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
```

测试`zh2en`功能：

![image-20250718211554588](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250718211554677.png)

发送的`json`数据：

```json
{
    "text":"基于HTTP/2构建，有HPACK压缩，支持单向、双向的流式传输，但不能直接调用浏览器的grpc服务，因此grpc常见于分布式系统的微服务中（实现服务器和服务器之间的交流）"
}
```

返回的`json`数据：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "original_text": "基于HTTP/2构建，有HPACK压缩，支持单向、双向的流式传输，但不能直接调用浏览器的grpc服务，因此grpc常见于分布式系统的微服务中（实现服务器和服务器之间的交流）",
        "translated_text": "Built on HTTP/2 with HPACK compression, it supports both unidirectional and bidirectional streaming. However, it cannot directly invoke gRPC services in browsers, so gRPC is commonly used in the microservices of distributed systems (enabling communication between servers).  \n\n(Note: The translation has been slightly refined for conciseness and clarity while maintaining technical accuracy.)",
        "from": "zh",
        "to": "en"
    }
}
```

测试`en2zh`功能：

![image-20250718211832919](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250718211832985.png)

发送的`json`数据：

```json
{
    "text":"The prices listed below are in unites of per 1M tokens. A token, the smallest unit of text that the model recognizes, can be a word, a number, or even a punctuation mark. We will bill based on the total number of input and output tokens by the model."
}
```

返回的`json`数据：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "original_text": "The prices listed below are in unites of per 1M tokens. A token, the smallest unit of text that the model recognizes, can be a word, a number, or even a punctuation mark. We will bill based on the total number of input and output tokens by the model.",
        "translated_text": "以下列出的价格单位为每百万个标记。标记是模型识别文本的最小单位，可以是一个词、一个数字，甚至是一个标点符号。我们将根据模型处理的输入与输出标记总数进行计费。",
        "from": "en",
        "to": "zh"
    }
}
```

## 3. 总结文本信息API

核心业务代码：

```go
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
```

其中，需要的数据结构为：

```go
// 总结请求
type SummaryRequest struct {
	Text string `json:"text" binding:"required"`
}

type SummaryResponse struct {
	OriginalText string `json:"original_text"`
	Summary      string `json:"summary"`
	Length       int    `json:"length"`
}
```

![image-20250718212653291](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250718212653390.png)

发送的`json`数据为：

```json
{
    "text":"麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。"
}
```

返回的`json`数据为：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "original_text": "麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。",
        "summary": "**总结：**  \n\n麦克白虽已掌权，却因预言中“班柯的子孙将统治国家”而心生恐惧，决定铲除班柯及其子弗里恩斯。尽管刺客杀死班柯，但弗里恩斯逃脱，令麦克白更加不安。在王室宴会上，班柯的鬼魂现身（仅麦克白可见），导致他当众失控，宴会被迫中断。  \n\n为寻求答案，麦克白再次求助女巫，得到三条预言：1）警惕麦克德夫；2）凡女人所生者皆无法伤害他；3）除非勃南森林移动，否则他永不败亡。这些看似绝对安全的预言让麦克白放松警惕，但女巫随后展示班柯后代世代为王的幻象，暗示麦克白的统治终将覆灭。  \n\n**核心要点：**  \n- 麦克白因预言对班柯家族赶尽杀绝，却未能消除隐患。  \n- 班柯鬼魂的出现暴露麦克白的精神崩溃，动摇其权威。  \n- 女巫的新预言表面安抚实则暗藏杀机，预示麦克白注定失败，班柯血脉将延续王权。",
        "length": 989
    }
}
```

