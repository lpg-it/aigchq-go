package aigchq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestQwenModelConstants(t *testing.T) {
	cases := map[string]string{
		ModelQwen37Plus:               "qwen3.7-plus",
		ModelQwen36Plus:               "qwen3.6-plus",
		ModelQwen35Flash:              "qwen3.5-flash",
		ModelQwen37PlusThinking:       "qwen3.7-plus-thinking",
		ModelQwen37PlusSearch:         "qwen3.7-plus-search",
		ModelQwen37PlusThinkingSearch: "qwen3.7-plus-thinking-search",
		ModelQwen37PlusFast:           "qwen3.7-plus-fast",
		ModelQwen37PlusFastSearch:     "qwen3.7-plus-fast-search",
		ProviderQwen:                  "qwen-web",
	}
	for got, want := range cases {
		if got != want {
			t.Fatalf("constant = %q, want %q", got, want)
		}
	}
}

func TestCreateQwenChatCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/qwen-web/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.Model != ModelQwen37PlusThinking {
			t.Fatalf("unexpected model: %q", req.Model)
		}
		if len(req.Messages) != 1 {
			t.Fatalf("unexpected messages: %+v", req.Messages)
		}
		writeJSON(w, ChatCompletionResponse{
			ID:       "chatcmpl_qwen",
			Object:   "chat.completion",
			Created:  time.Now().Unix(),
			Model:    req.Model,
			Provider: ProviderQwen,
			Choices: []Choice{{
				Index: 0,
				Message: Message{
					Role:             "assistant",
					Content:          "video summary",
					ReasoningContent: "watched key frames",
				},
				FinishReason: "stop",
			}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	resp, err := client.CreateQwenChatCompletion(context.Background(), &ChatCompletionRequest{
		Model: ModelQwen37PlusThinking,
		Messages: []Message{{
			Role: "user",
			Content: []ContentPart{
				{Type: "text", Text: "Analyze this video"},
				{
					Type: "input_file",
					File: &InputFile{
						Data:     "data:video/mp4;base64,AAAA",
						Filename: "clip.mp4",
						MimeType: "video/mp4",
					},
				},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Provider != ProviderQwen {
		t.Fatalf("provider = %q, want %q", resp.Provider, ProviderQwen)
	}
	if got := resp.Choices[0].Message.Content; got != "video summary" {
		t.Fatalf("content = %v, want video summary", got)
	}
	if got := resp.Choices[0].Message.ReasoningContent; got != "watched key frames" {
		t.Fatalf("reasoning = %q, want watched key frames", got)
	}
}

func TestListQwenModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/qwen-web/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		writeJSON(w, ModelsResponse{
			Object: "list",
			Data: []Model{{
				ID:       ModelQwen37Plus,
				Object:   "model",
				OwnedBy:  "qwen",
				Provider: ProviderQwen,
			}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	models, err := client.ListQwenModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(models.Data) != 1 || models.Data[0].ID != ModelQwen37Plus {
		t.Fatalf("unexpected models: %+v", models)
	}
}
