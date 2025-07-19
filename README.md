# helixlife
written test for helixlife

项目Github：[https://github.com/neKoui1/helixlife?tab=readme-ov-file#helixlife](https://github.com/neKoui1/helixlife?tab=readme-ov-file#helixlife)

在项目根目录下的配置文件`config.yaml`，读取使用的框架是`viper`，后端框架`gin`

因为要做到任务异步处理，需要对任务进行ID标识，使用的包为`uuid`

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

## 4. 实现翻译总结的异步任务

核心业务代码：

```go
// task_controller.go
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
```

```go
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
```

在调用的Update方法中，具体实现为：

```go
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
	}

}
```

![image-20250719153108692](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719153108771.png)

```json
{
    "type":"translate",
    "request": {
      "text": "麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。",
      "from": "zh",
      "to": "en"
    }
}
```

返回值：

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "b0d4179e-f8f3-4776-bda6-9d9694dfc515",
        "message": "任务已创建，请使用task_id查询进度",
        "status": "pending"
    }
}
```

![image-20250719153147001](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719153147091.png)

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "b0d4179e-f8f3-4776-bda6-9d9694dfc515",
        "type": "translate",
        "status": "completed",
        "response": {
            "original_text": "麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。",
            "translated_text": "Here is the English translation of the provided Chinese text:  \n\n---\n\nThough Macbeth had succeeded, he remained unsettled by the prophecy that \"Banquo's descendants would rule the kingdom.\" He thus concluded that Banquo posed a threat to his reign. Growing increasingly paranoid and ruthless, Macbeth invited Banquo to a royal banquet while secretly dispatching assassins to kill both him and his son Fleance. However, during the attack, Fleance managed to escape. Enraged, Macbeth grew deeply fearful of the potential threat posed by Banquo's lineage.  \n\nAt the banquet, Macbeth entertained the nobles with feasting and revelry. Unexpectedly, Banquo's ghost appeared and took Macbeth's seat. Since Macbeth was the only one who could see the specter, he was horrified and immediately descended into madness. His delirious ravings shocked and bewildered the assembled guests—they watched as he raged at an empty chair until Lady Macbeth intervened, assuring them, \"This is but a momentary fit; pay it no mind.\" Once Macbeth regained his composure, the feast resumed. But to everyone's dismay, the ghost reappeared, causing even greater chaos. Lady Macbeth had no choice but to dismiss the guests, ending the banquet in disgrace.  \n\nDeeply troubled, Macbeth sought out the witches once more, demanding the truth of their prophecies. To appease his fears, they summoned terrifying apparitions. First, an armored head warned him to \"beware Macduff.\" Next, a bloody child declared that \"none of woman born shall harm Macbeth.\" Finally, a crowned child assured him that he would \"never be vanquished until Birnam Wood comes to Dunsinane Hill.\" Macbeth took comfort in these omens, reasoning that all men are born of women and forests cannot move of their own accord.  \n\nWhen Macbeth inquired about Banquo's descendants, the witches conjured a procession of eight crowned kings, all resembling Banquo, with the last holding a mirror reflecting countless more. Macbeth recognized them as Banquo's heirs, ruling over vast and numerous realms.  \n\n---\n\nThis translation maintains the original meaning while ensuring natural flow and readability in English. Let me know if you'd like any refinements!",
            "from": "zh",
            "to": "en"
        },
        "progress": 100,
        "created_at": "2025-07-19T15:30:44.5325747+08:00",
        "updated_at": "2025-07-19T15:31:07.6339058+08:00",
        "started_at": "2025-07-19T15:30:44.5584566+08:00",
        "ended_at": "2025-07-19T15:31:07.6339058+08:00"
    }
}
```

Summary功能的测试：

![image-20250719153306581](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719153306656.png)

```json
{
    "type":"summary",
    "request": {
      "text": "麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。"
    }
}
```

![image-20250719153345292](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719153345386.png)

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "775eb5e3-0d4d-4d53-81a1-f4bb2d25eb12",
        "type": "summary",
        "status": "completed",
        "response": {
            "original_text": "麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。",
            "summary": "**总结：**\n\n麦克白因预言担忧班柯后裔威胁其统治，遂设计杀害班柯及其子弗里恩斯，但弗里恩斯逃脱。宴会上，麦克白因独见班柯鬼魂而失控，导致宴会混乱中断。为寻求安心，他再次求助女巫，获得三条预言：警惕麦克德夫、凡妇人所生者皆不能伤他、除非勃南森林移动否则不败。这些预言让麦克白暂时放松，但女巫随后展示班柯后裔世代为王的幻象，暗示其统治终将覆灭。麦克白虽表面释然，内心仍被恐惧与不安笼罩。",
            "length": 567
        },
        "progress": 100,
        "created_at": "2025-07-19T15:32:52.306825+08:00",
        "updated_at": "2025-07-19T15:33:02.8739229+08:00",
        "started_at": "2025-07-19T15:32:52.306825+08:00",
        "ended_at": "2025-07-19T15:33:02.8739229+08:00"
    }
}
```

## 5. 流式处理

ServerSentEvents实现流式响应

```go
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

```

之后在`ai_service`层中编写相关的业务处理

```go
// TranslateStream 流式翻译
func (s *AIService) TranslateStream(req *models.StreamRequest, callback func(string) error) error {
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
		Stream:      true,
	}
	resp, err := s.Client.ChatStream(*dsReq)
	if err != nil {
		return fmt.Errorf("流式翻译失败: %v", err)
	}
	return s.Client.ProcessStream(resp, func(data string) error {
		// 解析流式数据
		var chunk models.StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("解析流式数据失败: %v", err)
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			return callback(chunk.Choices[0].Delta.Content)
		}
		return nil
	})
}

// SummaryStream 流式总结
func (s *AIService) SummaryStream(req *models.StreamRequest, callback func(string) error) error {
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
		MaxTokens:   500,
		Stream:      true,
	}
	resp, err := s.Client.ChatStream(*dsReq)
	if err != nil {
		return fmt.Errorf("流式总结失败: %v", err)
	}
	return s.Client.ProcessStream(resp, func(data string) error {
		// 解析流式数据
		var chunk models.StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return fmt.Errorf("解析流式数据失败: %v", err)
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			return callback(chunk.Choices[0].Delta.Content)
		}
		return nil
	})
}
```

在`stream_controller`层中，使用channel不断接收流式返回的数据并输出：

```go
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
```

翻译发送的json请求：

```go
{
    "text":"麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。",
    "from":"zh",
    "to":"en"
}
```

返回结果为：

![image-20250719162400290](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719162400458.png)

总结发送的json请求和返回结果：

```json
{
    "text":"麦克白虽然成功了，却因“班柯的子孙将要君临一国”的预言仍感不安。他遂认为班柯会危害自己的统治。越发偏执冷酷的麦克白于是邀班柯席王室宴会，并暗地派刺客杀死他与儿子弗里恩斯；但是在刺杀中弗里恩斯却侥幸逃脱。麦克白倍感愤怒，深恐班柯后裔的潜在威胁。宴会中，麦克白宴请贵族吃喝做乐。没想到，班柯的鬼魂出现，并坐在麦克白的座位上。因为麦克白是唯一能看见那鬼魂的人，他受惊吓而立刻发疯。他的胡言乱语令在场的所有人倍感震惊、不解——所有人看着他对着空座椅咆哮，直到夫人告诉大家：“麦克白不过是偶然疾病，没什么大碍。”麦克白镇定之后，大家继续吃喝。却没想到鬼魂再次出现，导致更大的混乱。麦克白夫人只得请贵客离去，宴会尴尬收场。深受困扰的麦克白再次造访女巫们，也从她们那里获得预言的真相。女巫招出恐怖的幽灵以缓和麦克白的恐惧。首先，女巫招出披甲的头颅，提醒他要“留心麦克德夫”。然后，一个血淋淋的儿童告诉他“没有一个妇人所生下的人可以伤害麦克白”。第三，一个加冕的小孩告诉他“永远不会被打败，除非勃南的树林向邓西嫩的高山移动”。麦克白感到释然，因为所有人都是从妇人而生，而森林不会自己移动。当麦克白问及班戈的后裔时，女巫召唤了八个带着冠冕的国王，长相都与班戈相似，最后一个则拿着镜子，折射出更多的国王。麦克白认出这都是班柯的后裔，并且他们的国度繁多。"
}
```

![image-20250719162500668](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719162500782.png)

# Dify

一个自带GUI的开源LLM应用开发平台，省去了后端服务的开发过程

## 1. API文件夹中的内容

首先是dify/api/README.md中的文档，主要是教用户如何启动后端的web服务：配置环境变量，下载依赖等

![image-20250719165614604](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719165614668.png)

结合其他的文件夹内容（实际上看文件夹名就不难猜出来），web文件夹存放前端代码，sdks存放对接其他语言client的sdk，可以得知api文件夹中存放的是dify整个后端项目的核心代码，主要使用的语言是Python及其后端flask框架

按照一般的后端开发架构思路，开发顺序为model-service-controller，那么直接找到controller层查看具体的各种http请求响应

```python
class LoginApi(Resource):
    """Resource for web app email/password login."""

    @setup_required
    @only_edition_enterprise
    def post(self):
        """Authenticate user and login."""
        parser = reqparse.RequestParser()
        parser.add_argument("email", type=email, required=True, location="json")
        parser.add_argument("password", type=valid_password, required=True, location="json")
        args = parser.parse_args()

        try:
            account = WebAppAuthService.authenticate(args["email"], args["password"])
        except services.errors.account.AccountLoginError:
            raise AccountBannedError()
        except services.errors.account.AccountPasswordError:
            raise EmailOrPasswordMismatchError()
        except services.errors.account.AccountNotFoundError:
            raise AccountNotFound()

        token = WebAppAuthService.login(account=account)
        return {"result": "success", "data": {"access_token": token}}
```

在API文件夹中

* controllers/web: 实现了最终用户的用户登录、邮箱验证码、找回密码、信息获取和反馈、生成类似信息等控制器
* controllers/console: 实现了类似管理员的管理界面相关api，如应用管理（聊天助手、生成应用等）、workspace管理、用户管理和数据集管理等
* service: 实现了内部应用服务，应用生成服务的核心逻辑的函数为：

```python
class AppGenerateService:
    system_rate_limiter = RateLimiter("app_daily_rate_limiter", dify_config.APP_DAILY_RATE_LIMIT, 86400)

    @classmethod
    def generate(
        cls,
        app_model: App,
        user: Union[Account, EndUser],
        args: Mapping[str, Any],
        invoke_from: InvokeFrom,
        streaming: bool = True,
    ):
```

1. 付费用户检查
2. 根据app_model.mode调用不同的generate
   1. completion: 文本补全应用
   2. agent chat: 智能体对话
   3. chat: 对话型应用
   4. advanced_chat: 高级对话应用
   5. workflow: 工作流应用

在core/rag里面支持了强大的RAG功能：

```python
class DatasetRetrieval:
    def __init__(self, application_generate_entity=None):
        self.application_generate_entity = application_generate_entity

    def retrieve(
        self,
        app_id: str,
        user_id: str,
        tenant_id: str,
        model_config: ModelConfigWithCredentialsEntity,
        config: DatasetEntity,
        query: str,
        invoke_from: InvokeFrom,
        show_retrieve_source: bool,
        hit_callback: DatasetIndexToolCallbackHandler,
        message_id: str,
        memory: Optional[TokenBufferMemory] = None,
        inputs: Optional[Mapping[str, Any]] = None,
    ) -> Optional[str]:
```

支持single和multiple retrieve对数据集进行检索，之后对内外部的文档进行分别处理

在workflow_entry.py中，该项目同样支持工作流的事件执行和事件回调

```python
class WorkflowEntry:
    def __init__(
        self,
        tenant_id: str,
        app_id: str,
        workflow_id: str,
        workflow_type: WorkflowType,
        graph_config: Mapping[str, Any],
        graph: Graph,
        user_id: str,
        user_from: UserFrom,
        invoke_from: InvokeFrom,
        call_depth: int,
        variable_pool: VariablePool,
        thread_pool_id: Optional[str] = None,
    ) -> None:
    def run(
        self,
        *,
        callbacks: Sequence[WorkflowCallback],
    ) -> Generator[GraphEngineEvent, None, None]:
        """
        :param callbacks: workflow callbacks
        """
```

同样，在企业级架构中，很重要的一点是多租户架构的用户隔离以及权限控制关系，在model/account中

```python
class TenantAccountRole(enum.StrEnum):
    OWNER = "owner"
    ADMIN = "admin"
    EDITOR = "editor"
    NORMAL = "normal"
    DATASET_OPERATOR = "dataset_operator"

    @staticmethod
    def is_valid_role(role: str) -> bool:
        if not role:
            return False
        return role in {
            TenantAccountRole.OWNER,
            TenantAccountRole.ADMIN,
            TenantAccountRole.EDITOR,
            TenantAccountRole.NORMAL,
            TenantAccountRole.DATASET_OPERATOR,
        }

```

综上，该开源项目在api层中实现了一个基于python flask的LLM应用开发的后端平台，集成了不同的AI应用开发和强大的数据集/文件系统，自带RAG和对话能力和完善的用户安全系统。

## 2. 用户登录鉴权

这个项目有两个Login

1. 终点用户的login api
2. 管理者/开发人员的login api

在console/auth/login.python中：

![image-20250719202223511](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719202223627.png)

```python
class LoginApi(Resource):
    """Resource for user login."""

    @setup_required
    @email_password_login_enabled
    def post(self):
        """Authenticate user and login."""
        parser = reqparse.RequestParser()
        parser.add_argument("email", type=email, required=True, location="json")
        parser.add_argument("password", type=valid_password, required=True, location="json")
        parser.add_argument("remember_me", type=bool, required=False, default=False, location="json")
        parser.add_argument("invite_token", type=str, required=False, default=None, location="json")
        parser.add_argument("language", type=str, required=False, default="en-US", location="json")
        args = parser.parse_args()

        if dify_config.BILLING_ENABLED and BillingService.is_email_in_freeze(args["email"]):
            raise AccountInFreezeError()

        is_login_error_rate_limit = AccountService.is_login_error_rate_limit(args["email"])
        if is_login_error_rate_limit:
            raise EmailPasswordLoginLimitError()

        invitation = args["invite_token"]
        if invitation:
            invitation = RegisterService.get_invitation_if_token_valid(None, args["email"], invitation)

        if args["language"] is not None and args["language"] == "zh-Hans":
            language = "zh-Hans"
        else:
            language = "en-US"

        try:
            if invitation:
                data = invitation.get("data", {})
                invitee_email = data.get("email") if data else None
                if invitee_email != args["email"]:
                    raise InvalidEmailError()
                account = AccountService.authenticate(args["email"], args["password"], args["invite_token"])
            else:
                account = AccountService.authenticate(args["email"], args["password"])
        except services.errors.account.AccountLoginError:
            raise AccountBannedError()
        except services.errors.account.AccountPasswordError:
            AccountService.add_login_error_rate_limit(args["email"])
            raise EmailOrPasswordMismatchError()
        except services.errors.account.AccountNotFoundError:
            if FeatureService.get_system_features().is_allow_register:
                token = AccountService.send_reset_password_email(email=args["email"], language=language)
                return {"result": "fail", "data": token, "code": "account_not_found"}
            else:
                raise AccountNotFound()
        # SELF_HOSTED only have one workspace
        tenants = TenantService.get_join_tenants(account)
        if len(tenants) == 0:
            system_features = FeatureService.get_system_features()

            if system_features.is_allow_create_workspace and not system_features.license.workspaces.is_available():
                raise WorkspacesLimitExceeded()
            else:
                return {
                    "result": "fail",
                    "data": "workspace not found, please contact system admin to invite you to join in a workspace",
                }

        token_pair = AccountService.login(account=account, ip_address=extract_remote_ip(request))
        AccountService.reset_login_error_rate_limit(args["email"])
        return {"result": "success", "data": token_pair.model_dump()}
```

使用的是AccountService.authenticate函数进行密码验证，如果是邀请进来的用户会采取加盐后hash的方法生成密码，支持记录登录IP：

![image-20250719201821543](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719201821652.png)

```python
    @staticmethod
    def authenticate(email: str, password: str, invite_token: Optional[str] = None) -> Account:
        """authenticate account with email and password"""

        account = db.session.query(Account).filter_by(email=email).first()
        if not account:
            raise AccountNotFoundError()

        if account.status == AccountStatus.BANNED.value:
            raise AccountLoginError("Account is banned.")

        if password and invite_token and account.password is None:
            # if invite_token is valid, set password and password_salt
            salt = secrets.token_bytes(16)
            base64_salt = base64.b64encode(salt).decode()
            password_hashed = hash_password(password, salt)
            base64_password_hashed = base64.b64encode(password_hashed).decode()
            account.password = base64_password_hashed
            account.password_salt = base64_salt

        if account.password is None or not compare_password(password, account.password, account.password_salt):
            raise AccountPasswordError("Invalid email or password.")

        if account.status == AccountStatus.PENDING.value:
            account.status = AccountStatus.ACTIVE.value
            account.initialized_at = datetime.now(UTC).replace(tzinfo=None)

        db.session.commit()

        return cast(Account, account)
```

之后生成的tokenPair：

![image-20250719202316675](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719202316760.png)

```python
    @staticmethod
    def login(account: Account, *, ip_address: Optional[str] = None) -> TokenPair:
        if ip_address:
            AccountService.update_login_info(account=account, ip_address=ip_address)

        if account.status == AccountStatus.PENDING.value:
            account.status = AccountStatus.ACTIVE.value
            db.session.commit()

        access_token = AccountService.get_account_jwt_token(account=account)
        refresh_token = _generate_refresh_token()

        AccountService._store_refresh_token(refresh_token, account.id)

        return TokenPair(access_token=access_token, refresh_token=refresh_token)
```

可以直接看到采取的是jwt双token的鉴权方式

![image-20250719202635652](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719202635719.png)

```python
    @staticmethod
    def get_account_jwt_token(account: Account) -> str:
        exp_dt = datetime.now(UTC) + timedelta(minutes=dify_config.ACCESS_TOKEN_EXPIRE_MINUTES)
        exp = int(exp_dt.timestamp())
        payload = {
            "user_id": account.id,
            "exp": exp,
            "iss": dify_config.EDITION,
            "sub": "Console API Passport",
        }

        token: str = PassportService().issue(payload)
        return token
```

```python
def _generate_refresh_token(length: int = 64):
    token = secrets.token_hex(length)
    return token
```

可以看到，token携带的payload为用户id，expire time等

双token是一种为了延长用户会话而不用重新登录的常用方法

而用户的login只使用了单token:

![image-20250719203404046](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719203404138.png)

```python
class LoginApi(Resource):
    """Resource for web app email/password login."""

    @setup_required
    @only_edition_enterprise
    def post(self):
        """Authenticate user and login."""
        parser = reqparse.RequestParser()
        parser.add_argument("email", type=email, required=True, location="json")
        parser.add_argument("password", type=valid_password, required=True, location="json")
        args = parser.parse_args()

        try:
            account = WebAppAuthService.authenticate(args["email"], args["password"])
        except services.errors.account.AccountLoginError:
            raise AccountBannedError()
        except services.errors.account.AccountPasswordError:
            raise EmailOrPasswordMismatchError()
        except services.errors.account.AccountNotFoundError:
            raise AccountNotFound()

        token = WebAppAuthService.login(account=account)
        return {"result": "success", "data": {"access_token": token}}
```

```python
    @classmethod
    def login(cls, account: Account) -> str:
        access_token = cls._get_account_jwt_token(account=account)

        return access_token
```

![image-20250719203446477](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719203446558.png)

```python
    @classmethod
    def _get_account_jwt_token(cls, account: Account) -> str:
        exp_dt = datetime.now(UTC) + timedelta(hours=dify_config.ACCESS_TOKEN_EXPIRE_MINUTES * 24)
        exp = int(exp_dt.timestamp())

        payload = {
            "sub": "Web API Passport",
            "user_id": account.id,
            "session_id": account.email,
            "token_source": "webapp_login_token",
            "auth_type": "internal",
            "exp": exp,
        }

        token: str = PassportService().issue(payload)
        return token
```

这里携带的padload就是用户id和用户email以及对用户的auth_type进行了区分（可能是因为console管理人员已经对role做了区别所以这里对终点用户做了区分），并没有携带refresh token

## 3. command指令创建新workspace

创建workspace只能由图中的几个权限控制

![image-20250719205319518](https://cdn.jsdelivr.net/gh/neKoui1/picgo_images/img/20250719205319573.png)

