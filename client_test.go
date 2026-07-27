package aigchq

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClientDefaultTimeout(t *testing.T) {
	client, err := NewClient("test-key")
	if err != nil {
		t.Fatal(err)
	}
	if client.httpClient == nil {
		t.Fatal("http client is nil")
	}
	if client.httpClient.Timeout != 10*time.Minute {
		t.Fatalf("timeout = %s, want 10m", client.httpClient.Timeout)
	}
}

func TestCreateChatCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected authorization: %q", got)
		}
		if got := r.Header.Get("x-api-key"); got != "test-key" {
			t.Fatalf("unexpected x-api-key: %q", got)
		}
		if got := r.Header.Get("X-AIGCHQ-SDK"); got != "" {
			t.Fatalf("sdk must not inject provider/client marker headers, got %q", got)
		}
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Model != "gpt-5-5-thinking" || len(req.Messages) != 1 {
			t.Fatalf("unexpected request: %+v", req)
		}
		writeJSON(w, ChatCompletionResponse{
			ID:      "chatcmpl_1",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []Choice{{Index: 0, Message: Message{Role: "assistant", Content: "ok"}, FinishReason: "stop"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	resp, err := client.CreateChatCompletion(context.Background(), &ChatCompletionRequest{
		Model:    "gpt-5-5-thinking",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.Choices[0].Message.Content; got != "ok" {
		t.Fatalf("unexpected content: %v", got)
	}
}

func TestAsyncChatCompletionAndWait(t *testing.T) {
	var polls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/async/chat/completions":
			var req ChatCompletionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatal(err)
			}
			if req.Stream {
				t.Fatal("async chat request must force stream=false")
			}
			writeJSON(w, AsyncTaskResponse{
				TaskID:         "task_1",
				Status:         TaskStatusPending,
				PollURL:        "/v1/tasks/task_1",
				ConversationID: "conv_1",
				InputMessageID: "msg_in_1",
				CreatedAt:      time.Now(),
			})
		case "/v1/tasks/task_1":
			polls++
			if polls == 1 {
				writeJSON(w, TaskResponse{TaskID: "task_1", Status: TaskStatusProcessing, Type: TaskTypeChat, ConversationID: "conv_1", InputMessageID: "msg_in_1"})
				return
			}
			result, _ := json.Marshal(ChatCompletionResponse{
				ID:      "chatcmpl_2",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "gpt-5-5-thinking",
				Choices: []Choice{{Index: 0, Message: Message{Role: "assistant", Content: "done"}}},
			})
			writeJSON(w, TaskResponse{
				TaskID:          "task_1",
				Status:          TaskStatusCompleted,
				Type:            TaskTypeChat,
				Result:          result,
				ConversationID:  "conv_1",
				InputMessageID:  "msg_in_1",
				OutputMessageID: "msg_out_1",
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	resp, err := client.CreateChatCompletionAndWait(
		context.Background(),
		&ChatCompletionRequest{
			Model:    "gpt-5-5-thinking",
			Messages: []Message{{Role: "user", Content: "hi"}},
			Stream:   true,
		},
		WithPollInterval(time.Millisecond),
		WithPollTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.Choices[0].Message.Content; got != "done" {
		t.Fatalf("unexpected content: %v", got)
	}
	if polls != 2 {
		t.Fatalf("unexpected poll count: %d", polls)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeJSON(w, ErrorResponse{Error: ErrorDetail{Message: "rate limited", Type: "rate_limit_error", Code: "rate_limit"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	_, err := client.ListModels(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
	if apiErr.Detail == nil || apiErr.Detail.Code != "rate_limit" {
		t.Fatalf("unexpected detail: %+v", apiErr.Detail)
	}
}

func TestChatCompletionStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "text/event-stream" {
			t.Fatalf("unexpected accept: %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, `data: {"id":"chunk_1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"reasoning_content":"think"}}]}`+"\n\n")
		_, _ = io.WriteString(w, `data: {"id":"chunk_1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"hel"}}]}`+"\n\n")
		_, _ = io.WriteString(w, `data: {"id":"chunk_1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"lo"}}]}`+"\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	stream, err := client.CreateChatCompletionStream(context.Background(), &ChatCompletionRequest{
		Model:    "gpt-5-5-instant",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var text, reasoning string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		text += chunk.Choices[0].Delta.Content
		reasoning += chunk.Choices[0].Delta.ReasoningContent
	}
	if text != "hello" {
		t.Fatalf("unexpected stream text: %q", text)
	}
	if reasoning != "think" {
		t.Fatalf("unexpected stream reasoning: %q", reasoning)
	}
}

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	client, err := NewClient("test-key", WithBaseURL(baseURL), WithNoRetry())
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
