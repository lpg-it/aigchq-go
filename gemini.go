package aigchq

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// GeminiGenerateRequest is the native Gemini generateContent request shape.
type GeminiGenerateRequest struct {
	Contents          []GeminiContent         `json:"contents"`
	SystemInstruction *GeminiContent          `json:"systemInstruction,omitempty"`
	Tools             []GeminiTool            `json:"tools,omitempty"`
	ToolConfig        *GeminiToolConfig       `json:"toolConfig,omitempty"`
	GenerationConfig  *GeminiGenerationConfig `json:"generationConfig,omitempty"`
	SafetySettings    []GeminiSafetySetting   `json:"safetySettings,omitempty"`
}

type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents text, thoughts, attachments, and function calls in the
// native Gemini API. Thought is true for visible extended-thinking parts.
type GeminiPart struct {
	Text             string                  `json:"text,omitempty"`
	Thought          bool                    `json:"thought,omitempty"`
	InlineData       *GeminiInlineData       `json:"inlineData,omitempty"`
	FileData         *GeminiFileData         `json:"fileData,omitempty"`
	FunctionCall     *GeminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *GeminiFunctionResponse `json:"functionResponse,omitempty"`
}

type GeminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type GeminiFileData struct {
	MimeType string `json:"mimeType"`
	FileURI  string `json:"fileUri"`
}

type GeminiFunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

type GeminiFunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type GeminiTool struct {
	FunctionDeclarations []GeminiFunctionDeclaration `json:"functionDeclarations,omitempty"`
}

type GeminiFunctionDeclaration struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type GeminiToolConfig struct {
	FunctionCallingConfig *GeminiFunctionCallingConfig `json:"functionCallingConfig,omitempty"`
}

type GeminiFunctionCallingConfig struct {
	Mode                 string   `json:"mode,omitempty"`
	AllowedFunctionNames []string `json:"allowedFunctionNames,omitempty"`
}

type GeminiGenerationConfig struct {
	Temperature      *float64              `json:"temperature,omitempty"`
	TopP             *float64              `json:"topP,omitempty"`
	TopK             *int                  `json:"topK,omitempty"`
	MaxOutputTokens  *int                  `json:"maxOutputTokens,omitempty"`
	StopSequences    []string              `json:"stopSequences,omitempty"`
	CandidateCount   *int                  `json:"candidateCount,omitempty"`
	ResponseMimeType string                `json:"responseMimeType,omitempty"`
	ThinkingConfig   *GeminiThinkingConfig `json:"thinkingConfig,omitempty"`
}

// GeminiThinkingConfig controls native Gemini extended thinking.
type GeminiThinkingConfig struct {
	IncludeThoughts *bool  `json:"includeThoughts,omitempty"`
	ThinkingBudget  *int   `json:"thinkingBudget,omitempty"`
	ThinkingLevel   string `json:"thinkingLevel,omitempty"`
}

type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiGenerateResponse struct {
	Candidates     []GeminiCandidate     `json:"candidates"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
	UsageMetadata  *GeminiUsageMetadata  `json:"usageMetadata,omitempty"`
}

type GeminiCandidate struct {
	Content       GeminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason,omitempty"`
	Index         int                  `json:"index"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

type GeminiSafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type GeminiPromptFeedback struct {
	BlockReason   string               `json:"blockReason,omitempty"`
	SafetyRatings []GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

type GeminiUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
	ThoughtsTokenCount   int `json:"thoughtsTokenCount,omitempty"`
}

type GeminiModelsResponse struct {
	Models []GeminiModel `json:"models"`
}

type GeminiModel struct {
	Name                       string   `json:"name"`
	Version                    string   `json:"version,omitempty"`
	DisplayName                string   `json:"displayName,omitempty"`
	Description                string   `json:"description,omitempty"`
	InputTokenLimit            int      `json:"inputTokenLimit,omitempty"`
	OutputTokenLimit           int      `json:"outputTokenLimit,omitempty"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods,omitempty"`
	Temperature                float64  `json:"temperature,omitempty"`
	TopP                       float64  `json:"topP,omitempty"`
	TopK                       int      `json:"topK,omitempty"`
}

// GeminiStreamChunk is one native streamGenerateContent SSE event.
type GeminiStreamChunk struct {
	Candidates     []GeminiCandidate     `json:"candidates,omitempty"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
	UsageMetadata  *GeminiUsageMetadata  `json:"usageMetadata,omitempty"`
}

// GeminiResponseError is returned when a native Gemini stream emits an error
// event after the HTTP request has already succeeded.
type GeminiResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (e *GeminiResponseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Status != "" {
		return fmt.Sprintf("aigchq: Gemini stream: %s: %s", e.Status, e.Message)
	}
	return "aigchq: Gemini stream: " + e.Message
}

// ListGeminiModels calls the native Gemini-compatible model listing endpoint.
func (c *Client) ListGeminiModels(ctx context.Context) (*GeminiModelsResponse, error) {
	var out GeminiModelsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/v1beta/models", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGeminiModel calls the native Gemini-compatible model detail endpoint.
func (c *Client) GetGeminiModel(ctx context.Context, model string) (*GeminiModel, error) {
	model = strings.TrimPrefix(strings.TrimSpace(model), "models/")
	var out GeminiModel
	if err := c.doJSON(ctx, httpMethodGet, "/v1beta/models/"+pathEscape(model), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GenerateGeminiContent calls the native Gemini generateContent endpoint.
func (c *Client) GenerateGeminiContent(ctx context.Context, model string, req *GeminiGenerateRequest) (*GeminiGenerateResponse, error) {
	if req == nil {
		req = &GeminiGenerateRequest{}
	}
	model = strings.TrimPrefix(strings.TrimSpace(model), "models/")
	var out GeminiGenerateResponse
	path := "/v1beta/models/" + pathEscape(model) + ":generateContent"
	if err := c.doJSON(ctx, httpMethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GenerateGeminiContentStream calls the native Gemini streamGenerateContent
// endpoint and returns an SSE reader.
func (c *Client) GenerateGeminiContentStream(ctx context.Context, model string, req *GeminiGenerateRequest) (*GeminiContentStream, error) {
	if req == nil {
		req = &GeminiGenerateRequest{}
	}
	model = strings.TrimPrefix(strings.TrimSpace(model), "models/")
	path := "/v1beta/models/" + pathEscape(model) + ":streamGenerateContent"
	httpReq, payload, err := c.newRequest(ctx, httpMethodPost, path, nil, req)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.do(httpReq, payload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiErrorFromResponse(resp)
	}
	return &GeminiContentStream{
		body:   resp.Body,
		reader: bufio.NewReader(resp.Body),
	}, nil
}

// GeminiContentStream reads native Gemini streamGenerateContent SSE events.
type GeminiContentStream struct {
	body   io.ReadCloser
	reader *bufio.Reader
	done   bool
}

// Recv returns the next native Gemini chunk. It returns io.EOF when the
// connection closes or when an optional [DONE] event is received.
func (s *GeminiContentStream) Recv() (*GeminiStreamChunk, error) {
	if s == nil {
		return nil, errors.New("stream is nil")
	}
	if s.done {
		return nil, io.EOF
	}
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			s.done = true
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if err != nil {
				s.done = true
				return nil, err
			}
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			if err != nil {
				s.done = true
				return nil, err
			}
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			s.done = true
			return nil, io.EOF
		}

		var envelope struct {
			GeminiStreamChunk
			Error *GeminiResponseError `json:"error,omitempty"`
		}
		if err := json.Unmarshal([]byte(data), &envelope); err != nil {
			return nil, err
		}
		if envelope.Error != nil {
			s.done = true
			return nil, envelope.Error
		}
		return &envelope.GeminiStreamChunk, nil
	}
}

func (s *GeminiContentStream) Close() error {
	if s == nil || s.body == nil {
		return nil
	}
	s.done = true
	return s.body.Close()
}
