# AIGCHQ Go SDK

Official Go SDK for the AIGCHQ API.

Repository name: `aigchq-go`  
Module path: `github.com/lpg-it/aigchq-go`  
Package name: `aigchq`

This SDK only calls your AIGCHQ platform API. It does not call ChatGPT, Gemini, or any upstream provider API directly.

## Install

```bash
go get github.com/lpg-it/aigchq-go
```

## Client

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
			{Role: "user", Content: "Hello"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
```

For local or self-hosted deployments:

```go
client, err := aigchq.NewClient(
	os.Getenv("AIGCHQ_API_KEY"),
	aigchq.WithBaseURL("https://aigchq.com"),
	aigchq.WithTimeout(5*time.Minute),
)
```

By default the SDK sends only normal JSON headers plus `Authorization: Bearer <key>` and `x-api-key: <key>`. Custom headers are opt-in through `WithHeader`.

## Sync Chat

`/v1/chat/completions` is OpenAI-compatible and synchronous.

```go
resp, err := client.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "今天是几号？"},
	},
})
```

## Streaming Chat

```go
stream, err := client.CreateChatCompletionStream(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-instant",
	Messages: []aigchq.Message{
		{Role: "user", Content: "逐字输出一句话。"},
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
	if len(chunk.Choices) > 0 {
		fmt.Print(chunk.Choices[0].Delta.Content)
	}
}
```

## Async Chat

Use async endpoints for long-running upstream calls, Cloudflare/Nginx timeout avoidance, or queue-based workflows.

### 1. Create Task

SDK:

```go
task, err := client.CreateAsyncChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gpt-5-5-thinking",
	Messages: []aigchq.Message{
		{Role: "user", Content: "写一份详细方案。"},
	},
})
```

Raw API response:

```json
{
  "task_id": "task_01HX...",
  "status": "pending",
  "poll_url": "/v1/tasks/task_01HX...",
  "created_at": "2026-05-12T10:00:00Z"
}
```

### 2. Poll Task

SDK:

```go
taskState, err := client.GetTask(ctx, task.TaskID)
```

Processing response:

```json
{
  "task_id": "task_01HX...",
  "provider": "chatgpt-web",
  "model": "gpt-5-5-thinking",
  "type": "chat",
  "status": "processing",
  "progress": 0.25,
  "created_at": "2026-05-12T10:00:00Z"
}
```

Completed response:

```json
{
  "task_id": "task_01HX...",
  "provider": "chatgpt-web",
  "model": "gpt-5-5-thinking",
  "type": "chat",
  "status": "completed",
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

Failed response:

```json
{
  "task_id": "task_01HX...",
  "type": "chat",
  "status": "failed",
  "error": "upstream quota exceeded",
  "created_at": "2026-05-12T10:00:00Z",
  "completed_at": "2026-05-12T10:00:05Z"
}
```

### 3. Wait Helper

```go
resp, err := client.CreateChatCompletionAndWait(
	ctx,
	&aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-thinking",
		Messages: []aigchq.Message{
			{Role: "user", Content: "写一份详细方案。"},
		},
	},
	aigchq.WithPollInterval(2*time.Second),
	aigchq.WithPollTimeout(30*time.Minute),
	aigchq.WithPollHook(func(task *aigchq.TaskResponse) {
		fmt.Println(task.Status)
	}),
)
```

## Images

Synchronous:

```go
resp, err := client.CreateImageGeneration(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张简洁的 AI 控制台产品图",
	Size:   "1024x1024",
})
```

Asynchronous:

```go
resp, err := client.CreateImageGenerationAndWait(ctx, &aigchq.ImageGenerationRequest{
	Model:  "gemini-3.1-pro",
	Prompt: "一张简洁的 AI 控制台产品图",
})
```

Async image completion returns a normal OpenAI-compatible image response in `result`:

```json
{
  "task_id": "task_01HY...",
  "type": "image",
  "status": "completed",
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

## Provider-Specific Calls

Use provider-specific routes when you want to force Gemini or ChatGPT.

```go
gemini := client.Provider(aigchq.ProviderGemini)

resp, err := gemini.CreateChatCompletion(ctx, &aigchq.ChatCompletionRequest{
	Model: "gemini-3.1-pro",
	Messages: []aigchq.Message{
		{Role: "user", Content: "解释负载均衡。"},
	},
})
```

Async provider-specific calls are also available:

```go
resp, err := client.Provider(aigchq.ProviderChatGPT).
	CreateChatCompletionAndWait(ctx, req)
```

## Continuous Conversations

AIGCHQ responses can include `conversation_id` and `conversation`. Pass those back on the next request to continue on the same platform conversation. The server handles provider account stickiness and upstream conversation mapping.

```go
first, err := client.CreateChatCompletion(ctx, req)
if err != nil {
	return err
}

nextReq := &aigchq.ChatCompletionRequest{
	Model:          first.Model,
	ConversationID: first.ConversationID,
	Conversation:   first.Conversation,
	Messages: []aigchq.Message{
		{Role: "user", Content: "我刚才问了什么？"},
	},
}
next, err := client.CreateChatCompletion(ctx, nextReq)
```

## Models and Health

```go
models, err := client.ListModels(ctx)
model, err := client.GetModel(ctx, "gpt-5-5-thinking")
health, err := client.Health(ctx)
```

## Account and Platform APIs

The SDK wraps the platform APIs used by the web app:

```go
user, err := client.Register(ctx, &aigchq.RegisterRequest{
	Email: "dev@example.com",
	Name: "Dev",
	Password: "secret",
})

credentials, err := client.ListProviderCredentials(ctx)
logs, err := client.ListProviderRequestLogs(ctx, 100)
stats, err := client.ListProviderRequestStats(ctx, 24)
```

Provider credential management:

```go
enabled := true
result, err := client.UpsertProviderCredential(ctx, &aigchq.ProviderCredentialRequest{
	Provider:    aigchq.ProviderGemini,
	AccountName: "dev@example.com",
	AuthPayload: json.RawMessage(`{"cookie":"..."}`),
	Enabled:    &enabled,
})
```

Image host management:

```go
isDefault := true
host, err := client.UpsertImageHost(ctx, &aigchq.ImageHostRequest{
	Provider:  aigchq.ImageHostFImage,
	Name:      "default",
	BaseURL:   "https://...",
	APIToken:  os.Getenv("FIMAGE_API_TOKEN"),
	IsDefault: &isDefault,
})
```

Auth capture:

```go
capture, err := client.CreateAuthCapture(ctx, &aigchq.AuthCaptureRequest{
	Provider: aigchq.ProviderChatGPT,
})
fmt.Println(capture.CaptureURL)
```

## Errors

Non-2xx responses return `*aigchq.APIError`.

```go
resp, err := client.CreateChatCompletion(ctx, req)
if err != nil {
	var apiErr *aigchq.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode)
		if apiErr.Detail != nil {
			fmt.Println(apiErr.Detail.Type, apiErr.Detail.Code, apiErr.Detail.Message)
		}
		return err
	}
	return err
}
_ = resp
```

Network errors and context cancellations are returned directly from Go's HTTP stack.

## Retry

The client retries transient network errors, HTTP 429, and 5xx responses by default.

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

Disable retries:

```go
client, err := aigchq.NewClient(os.Getenv("AIGCHQ_API_KEY"), aigchq.WithNoRetry())
```

## Advanced

For newly added server endpoints before a typed SDK method exists:

```go
var out map[string]any
err := client.DoJSON(ctx, http.MethodGet, "/v1/models", nil, nil, &out)
```

## Examples

Runnable examples live in:

- `examples/chat`
- `examples/async`
- `examples/streaming`
- `examples/image`
- `examples/admin`

Run:

```bash
AIGCHQ_API_KEY=gf_xxx go run ./examples/chat
```

## Tests

```bash
go test ./...
```
