package aigchq

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGeminiModelConstants(t *testing.T) {
	got := []string{
		ModelGemini35FlashLite,
		ModelGemini36Flash,
		ModelGemini31Pro,
		ModelGemini35FlashLiteExtended,
		ModelGemini36FlashExtended,
		ModelGemini31ProExtended,
	}
	want := []string{
		"gemini-3.5-flash-lite",
		"gemini-3.6-flash",
		"gemini-3.1-pro",
		"gemini-3.5-flash-lite-extended",
		"gemini-3.6-flash-extended",
		"gemini-3.1-pro-extended",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("model constants = %#v, want %#v", got, want)
	}
}

func TestGeminiOpenAIRequestSupportsThinkingAndAttachments(t *testing.T) {
	models := []string{
		ModelGemini35FlashLite,
		ModelGemini36Flash,
		ModelGemini31Pro,
	}

	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			request := ChatCompletionRequest{
				Model: model,
				Messages: []Message{{
					Role: "user",
					Content: []ContentPart{
						{Type: "text", Text: "分析附件"},
						{
							Type: "input_image",
							ImageURL: &ImageURL{
								URL:    "data:image/png;base64,aGVsbG8=",
								Detail: "auto",
							},
							Name: "frame.png",
						},
						{
							Type: "input_file",
							File: &InputFile{
								Data:     "data:video/mp4;base64,aGVsbG8=",
								Filename: "clip.mp4",
								MimeType: "video/mp4",
							},
						},
					},
				}},
				ReasoningEffort: ReasoningEffortHigh,
				Thinking: &ThinkingConfig{
					Type:            "enabled",
					BudgetTokens:    testIntPtr(4096),
					IncludeThoughts: testBoolPtr(true),
				},
			}

			encoded, err := json.Marshal(request)
			if err != nil {
				t.Fatal(err)
			}
			var decoded map[string]any
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			if got := decoded["model"]; got != model {
				t.Fatalf("model = %#v, want %q", got, model)
			}
			if got := decoded["reasoning_effort"]; got != "high" {
				t.Fatalf("reasoning_effort = %#v", got)
			}
			thinking, ok := decoded["thinking"].(map[string]any)
			if !ok {
				t.Fatalf("thinking = %#v", decoded["thinking"])
			}
			if thinking["type"] != "enabled" || thinking["budget_tokens"] != float64(4096) || thinking["include_thoughts"] != true {
				t.Fatalf("thinking = %#v", thinking)
			}

			messages := decoded["messages"].([]any)
			content := messages[0].(map[string]any)["content"].([]any)
			image := content[1].(map[string]any)
			if image["type"] != "input_image" || image["name"] != "frame.png" {
				t.Fatalf("image part = %#v", image)
			}
			file := content[2].(map[string]any)
			filePayload := file["file"].(map[string]any)
			if file["type"] != "input_file" ||
				filePayload["filename"] != "clip.mp4" ||
				filePayload["mime_type"] != "video/mp4" ||
				filePayload["data"] != "data:video/mp4;base64,aGVsbG8=" {
				t.Fatalf("file part = %#v", file)
			}
		})
	}
}

func TestThinkingWireForms(t *testing.T) {
	tests := []struct {
		name     string
		thinking any
		assert   func(*testing.T, any)
	}{
		{
			name:     "boolean",
			thinking: true,
			assert: func(t *testing.T, got any) {
				t.Helper()
				if got != true {
					t.Fatalf("thinking = %#v, want true", got)
				}
			},
		},
		{
			name:     "string",
			thinking: "high",
			assert: func(t *testing.T, got any) {
				t.Helper()
				if got != "high" {
					t.Fatalf("thinking = %#v, want high", got)
				}
			},
		},
		{
			name: "object",
			thinking: &ThinkingConfig{
				Type:         "enabled",
				BudgetTokens: testIntPtr(4096),
			},
			assert: func(t *testing.T, got any) {
				t.Helper()
				object, ok := got.(map[string]any)
				if !ok || object["type"] != "enabled" || object["budget_tokens"] != float64(4096) {
					t.Fatalf("thinking = %#v, want enabled object", got)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := json.Marshal(ChatCompletionRequest{
				Model:    ModelGemini36Flash,
				Messages: []Message{{Role: "user", Content: "分析"}},
				Thinking: test.thinking,
			})
			if err != nil {
				t.Fatal(err)
			}
			var decoded map[string]any
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatal(err)
			}
			test.assert(t, decoded["thinking"])
		})
	}
}

func TestOpenAIReasoningContentResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ChatCompletionResponse{
			ID:    "chatcmpl_gemini",
			Model: ModelGemini31Pro,
			Choices: []Choice{{
				Index: 0,
				Message: Message{
					Role:             "assistant",
					ReasoningContent: "可见思考摘要",
					Content:          "最终回答",
				},
			}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	response, err := client.CreateChatCompletion(context.Background(), &ChatCompletionRequest{
		Model:    ModelGemini31Pro,
		Messages: []Message{{Role: "user", Content: "深入分析"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	message := response.Choices[0].Message
	if message.ReasoningContent != "可见思考摘要" || message.Content != "最终回答" {
		t.Fatalf("message = %#v", message)
	}
}

func TestGenerateGeminiContentSupportsNativeAttachmentsThoughtsAndThinking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/models/gemini-3.1-pro:generateContent" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
		}

		var request GeminiGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		parts := request.Contents[0].Parts
		if len(parts) != 3 || parts[1].InlineData == nil || parts[2].FileData == nil {
			t.Fatalf("parts = %#v", parts)
		}
		if parts[1].InlineData.MimeType != "video/mp4" || parts[1].InlineData.Data != "AAAA" {
			t.Fatalf("inlineData = %#v", parts[1].InlineData)
		}
		if parts[2].FileData.MimeType != "application/pdf" || parts[2].FileData.FileURI != "https://example.com/brief.pdf" {
			t.Fatalf("fileData = %#v", parts[2].FileData)
		}
		thinking := request.GenerationConfig.ThinkingConfig
		if thinking == nil || thinking.ThinkingLevel != "HIGH" || thinking.IncludeThoughts == nil || !*thinking.IncludeThoughts {
			t.Fatalf("thinkingConfig = %#v", thinking)
		}

		writeJSON(w, GeminiGenerateResponse{
			Candidates: []GeminiCandidate{{
				Index: 0,
				Content: GeminiContent{
					Role: "model",
					Parts: []GeminiPart{
						{Text: "可见思考摘要", Thought: true},
						{Text: "最终回答"},
					},
				},
				FinishReason: "STOP",
			}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	response, err := client.GenerateGeminiContent(context.Background(), "models/"+ModelGemini31Pro, &GeminiGenerateRequest{
		Contents: []GeminiContent{{
			Role: "user",
			Parts: []GeminiPart{
				{Text: "分析这些附件"},
				{InlineData: &GeminiInlineData{MimeType: "video/mp4", Data: "AAAA"}},
				{FileData: &GeminiFileData{MimeType: "application/pdf", FileURI: "https://example.com/brief.pdf"}},
			},
		}},
		GenerationConfig: &GeminiGenerationConfig{
			ThinkingConfig: &GeminiThinkingConfig{
				ThinkingLevel:   "HIGH",
				IncludeThoughts: testBoolPtr(true),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	parts := response.Candidates[0].Content.Parts
	if len(parts) != 2 || !parts[0].Thought || parts[0].Text != "可见思考摘要" || parts[1].Thought || parts[1].Text != "最终回答" {
		t.Fatalf("response parts = %#v", parts)
	}
}

func TestGenerateGeminiContentStreamParsesThoughtParts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/models/gemini-3.6-flash:streamGenerateContent" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "text/event-stream" {
			t.Fatalf("accept = %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, `data: {"candidates":[{"content":{"role":"model","parts":[{"text":"思考摘要","thought":true}]},"index":0}]}`+"\n\n")
		_, _ = io.WriteString(w, `data: {"candidates":[{"content":{"role":"model","parts":[{"text":"回答"}]},"finishReason":"STOP","index":0}]}`+"\n\n")
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	stream, err := client.GenerateGeminiContentStream(context.Background(), ModelGemini36Flash, &GeminiGenerateRequest{
		Contents: []GeminiContent{{
			Role:  "user",
			Parts: []GeminiPart{{Text: "你好"}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	first, err := stream.Recv()
	if err != nil {
		t.Fatal(err)
	}
	firstPart := first.Candidates[0].Content.Parts[0]
	if !firstPart.Thought || firstPart.Text != "思考摘要" {
		t.Fatalf("first part = %#v", firstPart)
	}

	second, err := stream.Recv()
	if err != nil {
		t.Fatal(err)
	}
	secondPart := second.Candidates[0].Content.Parts[0]
	if secondPart.Thought || secondPart.Text != "回答" {
		t.Fatalf("second part = %#v", secondPart)
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("final Recv error = %v, want io.EOF", err)
	}
}

func TestGeminiNativeModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1beta/models":
			writeJSON(w, GeminiModelsResponse{Models: []GeminiModel{{
				Name:                       "models/" + ModelGemini36Flash,
				DisplayName:                "3.6 Flash",
				SupportedGenerationMethods: []string{"generateContent", "streamGenerateContent"},
			}}})
		case "/v1beta/models/gemini-3.6-flash":
			writeJSON(w, GeminiModel{
				Name:        "models/" + ModelGemini36Flash,
				DisplayName: "3.6 Flash",
			})
		default:
			t.Fatalf("path = %q", r.URL.Path)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	models, err := client.ListGeminiModels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(models.Models) != 1 || models.Models[0].Name != "models/"+ModelGemini36Flash {
		t.Fatalf("models = %#v", models.Models)
	}
	model, err := client.GetGeminiModel(context.Background(), "models/"+ModelGemini36Flash)
	if err != nil {
		t.Fatal(err)
	}
	if model.DisplayName != "3.6 Flash" {
		t.Fatalf("model = %#v", model)
	}
}

func TestGeminiContentStreamReturnsNativeErrorEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, `data: {"error":{"code":400,"message":"bad attachment","status":"INVALID_ARGUMENT"}}`+"\n\n")
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	stream, err := client.GenerateGeminiContentStream(context.Background(), ModelGemini35FlashLite, &GeminiGenerateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	_, err = stream.Recv()
	var geminiErr *GeminiResponseError
	if !errors.As(err, &geminiErr) {
		t.Fatalf("error = %T %v, want *GeminiResponseError", err, err)
	}
	if geminiErr.Code != 400 || geminiErr.Status != "INVALID_ARGUMENT" {
		t.Fatalf("Gemini error = %#v", geminiErr)
	}
}

func testIntPtr(value int) *int    { return &value }
func testBoolPtr(value bool) *bool { return &value }
