# AIGCHQ Go SDK

`aigchq-go` 是 AIGCHQ API 的官方 Go SDK。

- GitHub: `github.com/lpg-it/aigchq-go`
- Go module: `github.com/lpg-it/aigchq-go`
- Package: `aigchq`
- 默认 API Base URL: `https://aigchq.com`

这个 SDK 的职责只有一个：帮助 Go 开发者快速调用 AIGCHQ 平台 API。它不会直接请求 ChatGPT、Gemini 或任何上游提供商接口，也不会处理上游 Cookie、网页请求头、浏览器指纹等内容。

## 安装

```bash
go get github.com/lpg-it/aigchq-go
```

## 最小可用示例

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	aigchq "github.com/lpg-it/aigchq-go"
)

func main() {
	client, err := aigchq.NewClient(os.Getenv("AIGCHQ_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.CreateChatCompletion(context.Background(), &aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-thinking",
		Messages: []aigchq.Message{
			{Role: "user", Content: "你好，用一句话介绍 AIGCHQ。"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

运行：

```bash
AIGCHQ_API_KEY=gf_xxx go run ./examples/chat
```

## Client 配置

### 默认客户端

```go
client, err := aigchq.NewClient(os.Getenv("AIGCHQ_API_KEY"))
```

SDK 默认会发送：

- `Accept: application/json`
- `Content-Type: application/json`，仅当有请求体时发送
- `Authorization: Bearer <api_key>`，仅当 API Key 非空时发送
- `x-api-key: <api_key>`，仅当 API Key 非空时发送

SDK 默认不会注入 `User-Agent`、`X-AIGCHQ-SDK` 或任何 provider 相关请求头。

### 私有部署地址

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithBaseURL("https://your-domain.example.com"),
)
```

### 超时

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithTimeout(10*time.Minute),
)
```

SDK 默认 HTTP 超时是 10 分钟。同步接口会占用当前 HTTP 请求生命周期；长任务仍建议使用异步接口，避免 Nginx、Cloudflare 或调用方超时。

### 自定义 HTTP Client

```go
httpClient := &http.Client{
	Timeout: 10 * time.Minute,
}

client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithHTTPClient(httpClient),
)
```

### 自定义 Header

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithHeader("X-Request-Source", "my-service"),
)
```

### 重试

默认会重试瞬时网络错误、HTTP 429 和 5xx，默认最多重试 2 次。

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithRetry(aigchq.RetryConfig{
		MaxRetries: 3,
		MinDelay:   500 * time.Millisecond,
		MaxDelay:   5 * time.Second,
	}),
)
```

关闭重试：

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithNoRetry(),
)
```

### 更新 API Key

```go
client.SetAPIKey("gf_new_key")
```

## 如何选择接口

| 场景 | 推荐方法 | 说明 |
| --- | --- | --- |
| 要兼容 OpenAI `/v1/chat/completions` | `CreateChatCompletion` | 同步接口，调用方等待最终响应 |
| 聊天可能超过 60 秒 | `CreateAsyncChatCompletion` + `WaitChatCompletion` | 创建任务后轮询，最适合外部系统集成 |
| 需要边生成边显示 | `CreateChatCompletionStream` | SSE 流式读取 chunk |
| 使用 Gemini 原生 `generateContent` | `GenerateGeminiContent` | 支持原生 parts、thinkingConfig 和 thought 响应 |
| 使用 Gemini 原生流式接口 | `GenerateGeminiContentStream` | 读取 `streamGenerateContent` SSE |
| 生成图片且可能耗时 | `CreateAsyncImageGeneration` + `WaitImageGeneration` | 避免长连接超时 |
| 强制用 Gemini 或 ChatGPT | `client.Provider(...).Create...` | 走 provider 专属路由 |
| 管理账号、日志、统计 | `ListProviderCredentials` 等平台方法 | 封装后台管理 API |

## Gemini 三模型、扩展思考和附件

SDK 为 Gemini Web 当前三个模型提供常量：

```go
aigchq.ModelGemini35FlashLite // gemini-3.5-flash-lite
aigchq.ModelGemini36Flash     // gemini-3.6-flash
aigchq.ModelGemini31Pro       // gemini-3.1-pro
```

三个模型使用相同的附件结构，均可在对话中传图片、视频、音频、PDF
和普通文件。推荐选择基础模型并设置 `ReasoningEffort` 开启扩展思考：

```go
req := &aigchq.ChatCompletionRequest{
	Model: aigchq.ModelGemini36Flash,
	Messages: []aigchq.Message{
		{Role: "user", Content: "深入分析这个方案"},
	},
	ReasoningEffort: aigchq.ReasoningEffortHigh,
}
```

`medium`、`high`、`xhigh`、`max`、`extended` 会选择扩展思考；
`none`、`minimal`、`low` 或省略字段会使用基础模式。SDK 也保留
`ModelGemini35FlashLiteExtended`、`ModelGemini36FlashExtended` 和
`ModelGemini31ProExtended` 三个显式兼容常量。

API 还接受 boolean、string 或 object 形式的 `thinking`。object 形式可使用
`ThinkingConfig`：

```go
budget := 4096
req.Thinking = &aigchq.ThinkingConfig{
	Type:         "enabled",
	BudgetTokens: &budget,
}
```

同步响应中的可见思考摘要位于
`resp.Choices[0].Message.ReasoningContent`，流式响应则位于
`chunk.Choices[0].Delta.ReasoningContent`。

## Chat 同步接口

对应 HTTP API：

```http
POST /v1/chat/completions
```

SDK 方法：

```go
resp, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "system", Content: "你是一个简洁的助手。"},
		{Role: "user", Content: "今天是几号？"},
	},
})
if err != nil {
	return err
}
fmt.Println(resp.Choices[0].Message.Content)
```

常用字段：

```go
req := &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "写一个 Go HTTP 示例。"},
	},
	Temperature: ptrFloat64(0.7),
	MaxTokens:   ptrInt(1000),
	WebSearch:   true,
	Metadata: map[string]any{
		"trace_id": "trace_123",
	},
}
```

多模态消息：

```go
resp, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: aigchq.ModelGemini36Flash,
	Messages: []aigchq.Message{
		{
			Role: "user",
			Content: []aigchq.ContentPart{
				{Type: "text", Text: "识别视频字幕，并参考图片说明场景"},
				{
					Type: "input_image",
					Name: "reference.png",
					ImageURL: &aigchq.ImageURL{
						URL: "data:image/png;base64,BASE64_IMAGE_BYTES",
					},
				},
				{
					Type: "input_file",
					File: &aigchq.InputFile{
						Data:     "data:video/mp4;base64,BASE64_VIDEO_BYTES",
						Filename: "episode-1.mp4",
						MimeType: "video/mp4",
					},
				},
			},
		},
	},
})
```

`ImageURL.URL` 和 `InputFile.URL` 也可以使用无需鉴权的公网 HTTP(S) URL。
本地文件应先转换为 data URI；不要传本机路径或 `file://` URL。包含本地
data URI 的附件必须使用同步接口。同步 Chat 的 JSON 请求体上限为 150 MiB，
异步任务请求体上限为 4 MiB 且附件只接受公网 HTTP(S) URL。Gemini provider
最多接收 8 个附件，单个附件及全部附件解码后的总大小均以 100 MiB 为上限；
base64/data URI 会使 JSON 请求体膨胀，请同时预留编码开销。

工具调用：

```go
resp, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "查询北京天气"},
	},
	Tools: []aigchq.Tool{
		{
			Type: "function",
			Function: aigchq.FunctionDef{
				Name:        "get_weather",
				Description: "Get weather by city",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"city": map[string]any{"type": "string"},
					},
					"required": []string{"city"},
				},
			},
		},
	},
})
```

## Chat 流式接口

对应 HTTP API：

```http
POST /v1/chat/completions
```

SDK 会自动设置 `stream=true` 并读取 SSE。

```go
stream, err := client.CreateChatCompletionStream(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-instant",
	Messages: []aigchq.Message{
		{Role: "user", Content: "逐字输出一个短故事。"},
	},
})
if err != nil {
	return err
}
defer stream.Close()

for {
	chunk, err := stream.Recv()
	if err == io.EOF {
		break
	}
	if err != nil {
		return err
	}
	if len(chunk.Choices) == 0 {
		continue
	}
	fmt.Print(chunk.Choices[0].Delta.Content)
}
```

## Gemini 原生 v1beta 接口

SDK 对 `/v1beta/models`、`generateContent` 和 `streamGenerateContent`
提供 typed client。原生 `inlineData` 可传 base64 图片、音频或视频；
`fileData.fileUri` 可传无需额外凭据即可下载的公网 HTTP(S) 文件 URL。

列出和读取模型：

```go
models, err := client.ListGeminiModels(ctx)
model, err := client.GetGeminiModel(ctx, aigchq.ModelGemini36Flash)
```

原生附件与扩展思考：

```go
includeThoughts := true

resp, err := client.GenerateGeminiContent(ctx, aigchq.ModelGemini31Pro, &aigchq.GeminiGenerateRequest{
	Contents: []aigchq.GeminiContent{
		{
			Role: "user",
			Parts: []aigchq.GeminiPart{
				{Text: "识别视频内容，并对照 PDF 总结"},
				{InlineData: &aigchq.GeminiInlineData{
					MimeType: "video/mp4",
					Data:     "BASE64_VIDEO_BYTES",
				}},
				{FileData: &aigchq.GeminiFileData{
					MimeType: "application/pdf",
					FileURI:  "https://example.com/brief.pdf",
				}},
			},
		},
	},
	GenerationConfig: &aigchq.GeminiGenerationConfig{
		ThinkingConfig: &aigchq.GeminiThinkingConfig{
			ThinkingLevel:   "HIGH",
			IncludeThoughts: &includeThoughts,
		},
	},
})
if err != nil {
	return err
}

for _, part := range resp.Candidates[0].Content.Parts {
	if part.Thought {
		fmt.Println("thinking:", part.Text)
		continue
	}
	fmt.Println("answer:", part.Text)
}
```

原生流式响应：

```go
stream, err := client.GenerateGeminiContentStream(ctx, aigchq.ModelGemini36Flash, request)
if err != nil {
	return err
}
defer stream.Close()

for {
	chunk, err := stream.Recv()
	if err == io.EOF {
		break
	}
	if err != nil {
		return err
	}
	for _, candidate := range chunk.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Thought {
				fmt.Print("[thinking] ")
			}
			fmt.Print(part.Text)
		}
	}
}
```

## Chat 异步接口

异步接口不是 OpenAI 标准接口，它是 AIGCHQ 为长任务提供的任务接口。

任务会持久化在平台数据库中，并和平台会话消息关联：创建任务时会生成用户消息，任务完成后会生成助手消息。服务重启后已经完成或失败的任务仍可查询；重启时仍处于 `processing` 的任务会被标记为失败，调用方可以重新创建任务。

对应 HTTP API：

```http
POST /v1/async/chat/completions
GET  /v1/tasks/{task_id}
```

### 第一步：创建任务

```go
task, err := client.CreateAsyncChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "写一份详细的系统设计方案。"},
	},
})
if err != nil {
	return err
}

fmt.Println(task.TaskID)
fmt.Println(task.PollURL)
```

创建任务响应：

```json
{
  "task_id": "task_01HX...",
  "status": "pending",
  "poll_url": "/v1/tasks/task_01HX...",
  "conversation_id": "conv_01HX...",
  "input_message_id": "msg_in_...",
  "created_at": "2026-05-12T10:00:00Z"
}
```

`conversation_id` 是平台侧长期会话 ID。后续如果你希望继续同一个对话，把它放回 chat/image 请求的 `conversation_id` 字段即可；平台会继续粘住同一个上游账号和上游对话。`input_message_id` 是本次请求创建的用户消息 ID。

### 第二步：轮询任务

```go
state, err := client.GetTask(ctx, task.TaskID)
if err != nil {
	return err
}

fmt.Println(state.Status)
```

处理中响应：

```json
{
  "task_id": "task_01HX...",
  "provider": "chatgpt-web",
  "model": "gpt-5-5-thinking",
  "type": "chat",
  "status": "processing",
  "progress": 0.25,
  "conversation_id": "conv_01HX...",
  "input_message_id": "msg_in_...",
  "created_at": "2026-05-12T10:00:00Z"
}
```

成功响应：

```json
{
  "task_id": "task_01HX...",
  "provider": "chatgpt-web",
  "model": "gpt-5-5-thinking",
  "type": "chat",
  "status": "completed",
  "conversation_id": "conv_01HX...",
  "input_message_id": "msg_in_...",
  "output_message_id": "msg_out_...",
  "result": {
    "id": "chatcmpl_...",
    "object": "chat.completion",
    "created": 1778569000,
    "model": "gpt-5-5-thinking",
    "choices": [
      {
        "index": 0,
        "message": {
          "role": "assistant",
          "content": "..."
        },
        "finish_reason": "stop"
      }
    ]
  },
  "created_at": "2026-05-12T10:00:00Z",
  "completed_at": "2026-05-12T10:00:40Z"
}
```

失败响应：

```json
{
  "task_id": "task_01HX...",
  "provider": "chatgpt-web",
  "model": "gpt-5-5-thinking",
  "type": "chat",
  "status": "failed",
  "error": "upstream quota exceeded",
  "conversation_id": "conv_01HX...",
  "input_message_id": "msg_in_...",
  "created_at": "2026-05-12T10:00:00Z",
  "completed_at": "2026-05-12T10:00:05Z"
}
```

### 第三步：转换结果

```go
state, err := client.WaitTask(ctx, task.TaskID)
if err != nil {
	return err
}

chat, err := state.ChatResult()
if err != nil {
	return err
}

fmt.Println(chat.Choices[0].Message.Content)
```

### 一步完成：创建并等待

```go
resp, err := client.CreateChatCompletionAndWait(
	ctx,
	&aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-thinking",
		Messages: []aigchq.Message{
			{Role: "user", Content: "写一份详细的系统设计方案。"},
		},
	},
	aigchq.WithPollInterval(2*time.Second),
	aigchq.WithPollTimeout(30*time.Minute),
	aigchq.WithPollHook(func(task *aigchq.TaskResponse) {
		fmt.Println("task status:", task.Status)
	}),
)
```

## 连续对话

第一次请求不要自己构造上游对话 ID。第一次请求完成后，如果响应里有 `conversation_id` 或 `conversation`，后续请求原样传回。

异步接口创建任务时就会返回 `ConversationID`，不必等任务完成；你可以保存 `task.ConversationID`，下一次请求传入 `ConversationID: task.ConversationID`。

```go
first, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "记住：我的项目叫 AIGCHQ。"},
	},
})
if err != nil {
	return err
}

second, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model:          first.Model,
	ConversationID: first.ConversationID,
	Conversation:   first.Conversation,
	Messages: []aigchq.Message{
		{Role: "user", Content: "我的项目叫什么？"},
	},
})
```

服务端负责：

- 复用同一个平台会话
- 粘住第一次选择的 provider 账号
- 维护上游 conversation id、parent message id 等映射

API 调用和前端调用都可以使用同一套连续对话机制。

## Image 同步接口

对应 HTTP API：

```http
POST /v1/images/generations
POST /v1/images/generate
```

SDK 方法：

```go
resp, err := client.CreateImageGeneration(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张干净的 AIGCHQ API 控制台产品海报",
	Size:   "1024x1024",
})
if err != nil {
	return err
}

for _, item := range resp.Data {
	fmt.Println(item.URL)
}
```

图生图：

```go
resp, err := client.CreateImageGeneration(ctx, &aigchq.ImageGenerationRequest{
	Model:     "gemini-3.1-pro",
	Prompt:   "把这张图改成科技产品发布会风格",
	ImageURL: "https://example.com/input.png",
})
```

## Image 异步接口

对应 HTTP API：

```http
POST /v1/async/images/generations
POST /v1/async/images/generate
GET  /v1/tasks/{task_id}
```

创建任务：

```go
task, err := client.CreateAsyncImageGeneration(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张干净的 AIGCHQ API 控制台产品海报",
	Size:   "1024x1024",
})
```

等待结果：

```go
resp, err := client.WaitImageGeneration(ctx, task.TaskID)
if err != nil {
	return err
}

for _, item := range resp.Data {
	fmt.Println(item.URL)
}
```

一步完成：

```go
resp, err := client.CreateImageGenerationAndWait(
	ctx,
	&aigchq.ImageGenerationRequest{
		Model:  "gemini-3.1-pro",
		Prompt: "一张干净的 AIGCHQ API 控制台产品海报",
	},
	aigchq.WithPollInterval(3*time.Second),
)
```

完成后的任务响应：

```json
{
  "task_id": "task_01HY...",
  "type": "image",
  "status": "completed",
  "conversation_id": "conv_01HY...",
  "input_message_id": "msg_in_...",
  "output_message_id": "msg_out_...",
  "result": {
    "created": 1778569000,
    "data": [
      {
        "url": "https://..."
      }
    ],
    "provider": "gemini-web"
  }
}
```

## Provider 专属调用

当你想强制使用某个 provider 时，使用 `client.Provider(providerName)`。

内置常量：

```go
const (
	aigchq.ProviderChatGPT = "chatgpt-web"
	aigchq.ProviderGemini  = "gemini-web"
)
```

### Provider 模型列表

HTTP API：

```http
GET /api/{provider}/models
```

SDK：

```go
models, err := client.Provider(aigchq.ProviderGemini).ListModels(ctx)
```

### Provider Chat 同步

HTTP API：

```http
POST /api/{provider}/chat/completions
```

SDK：

```go
resp, err := client.Provider(aigchq.ProviderGemini).CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gemini-3.1-pro",
	Messages: []aigchq.Message{
		{Role: "user", Content: "解释负载均衡。"},
	},
})
```

### Provider Chat 流式

```go
stream, err := client.Provider(aigchq.ProviderChatGPT).CreateChatCompletionStream(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-instant",
	Messages: []aigchq.Message{
		{Role: "user", Content: "流式输出。"},
	},
})
```

### Provider Chat 异步

HTTP API：

```http
POST /api/{provider}/async/chat/completions
GET  /v1/tasks/{task_id}
```

SDK：

```go
resp, err := client.Provider(aigchq.ProviderChatGPT).CreateChatCompletionAndWait(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "写一个方案。"},
	},
})
```

### Provider Image 同步

HTTP API：

```http
POST /api/{provider}/images/generations
POST /api/{provider}/images/generate
```

SDK：

```go
resp, err := client.Provider(aigchq.ProviderGemini).CreateImageGeneration(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张产品图",
})
```

### Provider Image 异步

HTTP API：

```http
POST /api/{provider}/async/images/generations
POST /api/{provider}/async/images/generate
GET  /v1/tasks/{task_id}
```

SDK：

```go
resp, err := client.Provider(aigchq.ProviderGemini).CreateImageGenerationAndWait(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张产品图",
})
```

## Models

### 列出全部模型

HTTP API：

```http
GET /v1/models
```

SDK：

```go
models, err := client.ListModels(ctx)
if err != nil {
	return err
}
for _, model := range models.Data {
	fmt.Println(model.ID, model.Provider, model.Capabilities)
}
```

### 获取单个模型

HTTP API：

```http
GET /v1/models/{model}
```

SDK：

```go
model, err := client.GetModel(ctx, "gpt-5-5-thinking")
```

## Health

HTTP API：

```http
GET /health
```

SDK：

```go
health, err := client.Health(ctx)
fmt.Println(health.Status, health.Version, health.Timestamp)
```

## Task API

### 获取任务

HTTP API：

```http
GET /v1/tasks/{task_id}
```

SDK：

```go
task, err := client.GetTask(ctx, "task_01HX...")
```

### 等待任意任务完成

```go
task, err := client.WaitTask(
	ctx,
	"task_01HX...",
	aigchq.WithPollInterval(2*time.Second),
	aigchq.WithPollTimeout(30*time.Minute),
)
```

### 等待 Chat 任务

```go
chat, err := client.WaitChatCompletion(ctx, "task_01HX...")
```

### 等待 Image 任务

```go
image, err := client.WaitImageGeneration(ctx, "task_01HY...")
```

## Auth API

这些接口用于平台账号体系，不是 provider 账号授权。

### 注册

HTTP API：

```http
POST /api/auth/register
```

SDK：

```go
user, err := client.Register(ctx, &aigchq.RegisterRequest{
	Email:    "dev@example.com",
	Name:     "Dev",
	Password: "secret",
})
```

注册成功后，如果响应里包含完整 API Key，SDK 会自动调用 `SetAPIKey` 更新当前客户端。

### 登录

HTTP API：

```http
POST /api/auth/login
```

SDK：

```go
user, err := client.Login(ctx, &aigchq.LoginRequest{
	Email:    "dev@example.com",
	Password: "secret",
})
```

### 登出

HTTP API：

```http
POST /api/auth/logout
```

SDK：

```go
err := client.Logout(ctx)
```

### 当前用户

HTTP API：

```http
GET /api/auth/me
```

SDK：

```go
user, err := client.Me(ctx)
```

### 更新当前用户

HTTP API：

```http
PATCH /api/auth/me
```

SDK：

```go
user, err := client.UpdateMe(ctx, &aigchq.UpdateUserRequest{
	Name: "New Name",
})
```

### 重置 API Key

HTTP API：

```http
POST /api/auth/api-key
```

SDK：

```go
user, err := client.RegenerateAPIKey(ctx)
fmt.Println(user.APIKey)
```

重置成功后，如果响应里包含完整 API Key，SDK 会自动调用 `SetAPIKey` 更新当前客户端。

## Provider Credential API

这些接口用于管理你在 AIGCHQ 平台里接入的 ChatGPT/Gemini 账号。

### Provider 状态

HTTP API：

```http
GET /api/provider-status
```

SDK：

```go
items, err := client.ProviderStatus(ctx)
for _, item := range items {
	fmt.Println(item.Provider, item.Configured, item.Accounts, item.Capabilities)
}
```

### 账号列表

HTTP API：

```http
GET /api/provider-credentials
```

SDK：

```go
items, err := client.ListProviderCredentials(ctx)
for _, item := range items {
	fmt.Println(item.Provider, item.AccountName, item.AccountEmail, item.Plan, item.Enabled)
	if item.Health != nil {
		fmt.Println(item.Health.Status, item.Health.TotalRequests, item.Health.TotalSuccesses)
	}
}
```

### 新增或更新账号

HTTP API：

```http
POST /api/provider-credentials
```

SDK：

```go
enabled := true
result, err := client.UpsertProviderCredential(ctx, &aigchq.ProviderCredentialRequest{
	Provider:     aigchq.ProviderGemini,
	AccountName:  "dev@example.com",
	AccountEmail: "dev@example.com",
	Plan:         "pro",
	AuthPayload:  json.RawMessage(`{"auth_payload_from_extension": true}`),
	Enabled:      &enabled,
})
if err != nil {
	return err
}
fmt.Println(result.Item.ID)
```

`AuthPayload` 是 AIGCHQ Auth Helper 扩展采集后提交给平台的授权数据。SDK 只把该 payload 发给你的 AIGCHQ 服务端保存，不解析、不请求上游 provider。

### 更新账号设置

HTTP API：

```http
PATCH /api/provider-credentials/{id}/settings
```

SDK：

```go
enabled := true
result, err := client.UpdateProviderCredentialSettings(ctx, "credential_id", &aigchq.ProviderCredentialSettingsRequest{
	Enabled:    &enabled,
	ProxyURL:   "socks5://user:pass@host:port",
	ProxyScope: "all",
})
```

### 删除账号

HTTP API：

```http
DELETE /api/provider-credentials/{id}
```

SDK：

```go
err := client.DeleteProviderCredential(ctx, "credential_id")
```

## Request Log and Stats API

### 最近请求日志

HTTP API：

```http
GET /api/provider-request-logs?limit=100
```

SDK：

```go
logs, err := client.ListProviderRequestLogs(ctx, 100)
for _, logItem := range logs {
	fmt.Println(
		logItem.CreatedAt,
		logItem.Provider,
		logItem.AccountEmail,
		logItem.Model,
		logItem.RequestType,
		logItem.Status,
		logItem.DurationMS,
		logItem.ErrorMessage,
	)
}
```

### 请求统计

HTTP API：

```http
GET /api/provider-request-stats?hours=24
```

SDK：

```go
stats, err := client.ListProviderRequestStats(ctx, 24)
if err != nil {
	return err
}
for _, item := range stats.Items {
	successRate := 0.0
	if item.Total > 0 {
		successRate = float64(item.Success) / float64(item.Total)
	}
	fmt.Println(item.Provider, item.AccountEmail, item.Total, successRate)
}
```

## Image Host API

### 图片托管列表

HTTP API：

```http
GET /api/image-hosts
```

SDK：

```go
hosts, err := client.ListImageHosts(ctx)
```

### 新增或更新图片托管

HTTP API：

```http
POST /api/image-hosts
```

SDK：

```go
enabled := true
isDefault := true
result, err := client.UpsertImageHost(ctx, &aigchq.ImageHostRequest{
	Provider:  aigchq.ImageHostFImage,
	Name:      "default",
	BaseURL:   "https://your-fimage.example.com",
	APIToken:  os.Getenv("FIMAGE_API_TOKEN"),
	Enabled:   &enabled,
	IsDefault: &isDefault,
})
```

### 删除图片托管

HTTP API：

```http
DELETE /api/image-hosts/{id}
```

SDK：

```go
err := client.DeleteImageHost(ctx, "image_host_id")
```

## Conversation API

这些接口用于管理 AIGCHQ 平台自己的会话列表和消息列表。

### 会话列表

HTTP API：

```http
GET /api/conversations?limit=20&offset=0&pinned=false
```

SDK：

```go
pinned := false
resp, err := client.ListConversations(ctx, aigchq.ConversationListOptions{
	Limit:  20,
	Offset: 0,
	Pinned: &pinned,
})
for _, conv := range resp.Items {
	fmt.Println(conv.ID, conv.Title, conv.Model, conv.MessageCount)
}
```

### 会话消息

HTTP API：

```http
GET /api/conversations/{id}/messages
```

SDK：

```go
messages, err := client.ListConversationMessages(ctx, "conversation_id")
```

### 新增或更新会话

HTTP API：

```http
POST /api/conversations
```

SDK：

```go
conv, err := client.UpsertConversation(ctx, &aigchq.Conversation{
	ID:    "conversation_id",
	Title: "新的会话",
	Model: "gpt-5-5-thinking",
})
```

### 置顶或取消置顶

HTTP API：

```http
PATCH /api/conversations/{id}/pin
```

SDK：

```go
conv, err := client.PinConversation(ctx, "conversation_id", true)
```

### 删除会话

HTTP API：

```http
DELETE /api/conversations/{id}
```

SDK：

```go
err := client.DeleteConversation(ctx, "conversation_id")
```

## Auth Capture API

Auth Capture 用于创建一次 provider 授权采集流程。浏览器扩展采集完成后会回调 complete URL。

### 创建采集任务

HTTP API：

```http
POST /api/auth-captures
```

SDK：

```go
capture, err := client.CreateAuthCapture(ctx, &aigchq.AuthCaptureRequest{
	Provider: aigchq.ProviderChatGPT,
})
if err != nil {
	return err
}

fmt.Println(capture.Capture.ID)
fmt.Println(capture.CaptureURL)
fmt.Println(capture.CompleteURL)
```

### 查询采集状态

HTTP API：

```http
GET /api/auth-captures/{id}
```

SDK：

```go
capture, err := client.GetAuthCapture(ctx, "capture_id")
fmt.Println(capture.Status)
```

### 完成采集

HTTP API：

```http
POST /api/auth-captures/{id}/complete
```

SDK：

```go
result, err := client.CompleteAuthCapture(ctx, "capture_id", map[string]any{
	"provider": "chatgpt-web",
	"payload":  "...",
})
fmt.Println(result.OK, result.Provider, result.AccountName)
```

通常由 AIGCHQ Auth Helper 扩展调用该接口，业务后端不需要手动调用。

## 高级：直接调用新接口

如果服务端已经新增接口，但 SDK 还没发布对应 typed method，可以使用 `DoJSON`：

```go
var out map[string]any
err := client.DoJSON(ctx, http.MethodGet, "/v1/models", nil, nil, &out)
```

带 query：

```go
query := url.Values{}
query.Set("limit", "100")

var out map[string]any
err := client.DoJSON(ctx, http.MethodGet, "/api/provider-request-logs", query, nil, &out)
```

## 错误处理

非 2xx 响应会返回 `*aigchq.APIError`。

```go
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
	var apiErr *aigchq.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("status:", apiErr.StatusCode)
		fmt.Println("request id:", apiErr.RequestID)
		if apiErr.Detail != nil {
			fmt.Println("type:", apiErr.Detail.Type)
			fmt.Println("code:", apiErr.Detail.Code)
			fmt.Println("message:", apiErr.Detail.Message)
		}
		return err
	}
	return err
}
_ = resp
```

异步任务失败时，`WaitTask`、`WaitChatCompletion`、`WaitImageGeneration` 会返回任务里的 `error` 文本：

```go
resp, err := client.WaitChatCompletion(ctx, taskID)
if err != nil {
	return fmt.Errorf("async chat failed: %w", err)
}
_ = resp
```

网络错误、DNS 错误、TLS 错误、`context.Canceled`、`context.DeadlineExceeded` 会按 Go 标准错误返回。

## 完整方法索引

| SDK 方法 | HTTP API | 说明 |
| --- | --- | --- |
| `NewClient` | - | 创建客户端 |
| `MustNewClient` | - | 创建客户端，失败 panic |
| `SetAPIKey` | - | 更新当前客户端 API Key |
| `DoJSON` | 任意路径 | 调用未封装的新接口 |
| `Health` | `GET /health` | 健康检查 |
| `ListModels` | `GET /v1/models` | 全部可用模型 |
| `GetModel` | `GET /v1/models/{model}` | 单个模型 |
| `CreateChatCompletion` | `POST /v1/chat/completions` | 同步 chat |
| `CreateChatCompletionStream` | `POST /v1/chat/completions` | 流式 chat |
| `CreateAsyncChatCompletion` | `POST /v1/async/chat/completions` | 创建异步 chat 任务 |
| `CreateChatCompletionAndWait` | `POST /v1/async/chat/completions` + `GET /v1/tasks/{id}` | 创建并轮询 chat |
| `CreateImageGeneration` | `POST /v1/images/generations` | 同步图片生成 |
| `CreateAsyncImageGeneration` | `POST /v1/async/images/generations` | 创建异步图片任务 |
| `CreateImageGenerationAndWait` | `POST /v1/async/images/generations` + `GET /v1/tasks/{id}` | 创建并轮询图片 |
| `GetTask` | `GET /v1/tasks/{task_id}` | 查询任务 |
| `WaitTask` | `GET /v1/tasks/{task_id}` | 轮询任意任务 |
| `WaitChatCompletion` | `GET /v1/tasks/{task_id}` | 轮询 chat 任务并解析 |
| `WaitImageGeneration` | `GET /v1/tasks/{task_id}` | 轮询 image 任务并解析 |
| `Provider(name).ListModels` | `GET /api/{provider}/models` | provider 模型 |
| `Provider(name).CreateChatCompletion` | `POST /api/{provider}/chat/completions` | provider 同步 chat |
| `Provider(name).CreateChatCompletionStream` | `POST /api/{provider}/chat/completions` | provider 流式 chat |
| `Provider(name).CreateAsyncChatCompletion` | `POST /api/{provider}/async/chat/completions` | provider 异步 chat |
| `Provider(name).CreateChatCompletionAndWait` | provider async + task | provider 创建并轮询 chat |
| `Provider(name).CreateImageGeneration` | `POST /api/{provider}/images/generations` | provider 同步图片 |
| `Provider(name).CreateAsyncImageGeneration` | `POST /api/{provider}/async/images/generations` | provider 异步图片 |
| `Provider(name).CreateImageGenerationAndWait` | provider async + task | provider 创建并轮询图片 |
| `Register` | `POST /api/auth/register` | 注册平台用户 |
| `Login` | `POST /api/auth/login` | 登录平台用户 |
| `Logout` | `POST /api/auth/logout` | 登出 |
| `Me` | `GET /api/auth/me` | 当前用户 |
| `UpdateMe` | `PATCH /api/auth/me` | 更新当前用户 |
| `RegenerateAPIKey` | `POST /api/auth/api-key` | 重置 API Key |
| `ProviderStatus` | `GET /api/provider-status` | provider 配置状态 |
| `ListProviderCredentials` | `GET /api/provider-credentials` | provider 账号列表 |
| `UpsertProviderCredential` | `POST /api/provider-credentials` | 新增或更新 provider 账号 |
| `UpdateProviderCredentialSettings` | `PATCH /api/provider-credentials/{id}/settings` | 更新启用状态和代理 |
| `DeleteProviderCredential` | `DELETE /api/provider-credentials/{id}` | 删除 provider 账号 |
| `ListProviderRequestLogs` | `GET /api/provider-request-logs` | 最近请求日志 |
| `ListProviderRequestStats` | `GET /api/provider-request-stats` | 请求统计 |
| `ListImageHosts` | `GET /api/image-hosts` | 图片托管列表 |
| `UpsertImageHost` | `POST /api/image-hosts` | 新增或更新图片托管 |
| `DeleteImageHost` | `DELETE /api/image-hosts/{id}` | 删除图片托管 |
| `ListConversations` | `GET /api/conversations` | 会话列表 |
| `ListConversationMessages` | `GET /api/conversations/{id}/messages` | 会话消息 |
| `UpsertConversation` | `POST /api/conversations` | 新增或更新会话 |
| `PinConversation` | `PATCH /api/conversations/{id}/pin` | 置顶或取消置顶 |
| `DeleteConversation` | `DELETE /api/conversations/{id}` | 删除会话 |
| `CreateAuthCapture` | `POST /api/auth-captures` | 创建 provider 授权采集 |
| `GetAuthCapture` | `GET /api/auth-captures/{id}` | 查询授权采集 |
| `CompleteAuthCapture` | `POST /api/auth-captures/{id}/complete` | 完成授权采集 |

## 类型速查

### ChatCompletionRequest

常用字段：

- `Model`: 模型名
- `Messages`: OpenAI-compatible messages
- `Stream`: 是否流式，流式方法会自动设置
- `Temperature`, `TopP`, `MaxTokens`
- `Tools`, `ToolChoice`, `ParallelToolCalls`
- `ResponseFormat`
- `ReasoningEffort`, `Thinking`: 扩展思考控制
- `Image`, `ImageName`, `Images`, `Media`: 请求级多模态兼容字段
- `Provider`: 指定 provider，通常优先使用 `client.Provider(...)`
- `Account`: 指定账号名或邮箱
- `ConversationID`, `ParentMessageID`, `Conversation`: 连续对话
- `WebSearch`: 是否启用搜索
- `Metadata`, `ExtraBody`: 扩展字段

### ImageGenerationRequest

常用字段：

- `Model`
- `Prompt`
- `N`
- `Size`
- `Quality`
- `ResponseFormat`
- `Style`
- `Messages`
- `Width`, `Height`, `AspectRatio`, `Resolution`
- `ImageURL`, `ImageData`
- `Provider`, `Account`
- `ConversationID`, `ParentMessageID`, `Conversation`
- `DownloadMedia`

### TaskResponse

字段：

- `TaskID`
- `Provider`
- `Model`
- `Type`: `chat`、`image`、`video`
- `Status`: `pending`、`processing`、`completed`、`failed`
- `Progress`
- `Result`
- `Error`
- `ConversationID`: 平台侧长期会话 ID，可用于继续同一个对话
- `InputMessageID`: 本次任务对应的用户消息 ID
- `OutputMessageID`: 本次任务完成后生成的助手消息 ID
- `CreatedAt`
- `CompletedAt`

辅助方法：

```go
chat, err := task.ChatResult()
image, err := task.ImageResult()
done := task.IsTerminal()
```

## 示例目录

- `examples/chat`: 同步 chat
- `examples/async`: 异步 chat 创建和轮询
- `examples/streaming`: 流式 chat
- `examples/image`: 异步图片生成
- `examples/admin`: provider 账号和请求日志

运行：

```bash
AIGCHQ_API_KEY=gf_xxx go run ./examples/chat
AIGCHQ_API_KEY=gf_xxx go run ./examples/async
AIGCHQ_API_KEY=gf_xxx go run ./examples/streaming
AIGCHQ_API_KEY=gf_xxx go run ./examples/image
AIGCHQ_API_KEY=gf_xxx go run ./examples/admin
```

私有部署：

```bash
AIGCHQ_API_KEY=gf_xxx AIGCHQ_BASE_URL=https://your-domain.example.com go run ./examples/chat
```

## 测试

```bash
go test ./...
```

## 发布

```bash
git tag v1.0.0
git push origin main
git push origin v1.0.0
```
