package aigchq

import (
	"encoding/json"
	"time"
)

type Message struct {
	Role             string     `json:"role"`
	Content          any        `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	Name             string     `json:"name,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type ChatCompletionRequest struct {
	Model             string             `json:"model"`
	Messages          []Message          `json:"messages"`
	Stream            bool               `json:"stream,omitempty"`
	Temperature       *float64           `json:"temperature,omitempty"`
	TopP              *float64           `json:"top_p,omitempty"`
	N                 *int               `json:"n,omitempty"`
	MaxTokens         *int               `json:"max_tokens,omitempty"`
	PresencePenalty   *float64           `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64           `json:"frequency_penalty,omitempty"`
	LogitBias         map[string]float64 `json:"logit_bias,omitempty"`
	Stop              []string           `json:"stop,omitempty"`
	User              string             `json:"user,omitempty"`
	Tools             []Tool             `json:"tools,omitempty"`
	ToolChoice        any                `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool              `json:"parallel_tool_calls,omitempty"`
	ResponseFormat    *ResponseFormat    `json:"response_format,omitempty"`
	ReasoningEffort   string             `json:"reasoning_effort,omitempty"`
	Audio             map[string]any     `json:"audio,omitempty"`
	Modalities        []string           `json:"modalities,omitempty"`
	Provider          string             `json:"provider,omitempty"`
	Account           string             `json:"account,omitempty"`
	ConversationID    string             `json:"conversation_id,omitempty"`
	ParentMessageID   string             `json:"parent_message_id,omitempty"`
	Conversation      map[string]any     `json:"conversation,omitempty"`
	WebSearch         bool               `json:"web_search,omitempty"`
	Metadata          map[string]any     `json:"metadata,omitempty"`
	ExtraBody         map[string]any     `json:"extra_body,omitempty"`
	Image             string             `json:"image,omitempty"`
	ImageName         string             `json:"image_name,omitempty"`
	Images            []any              `json:"images,omitempty"`
	Media             []any              `json:"media,omitempty"`
	APIKey            any                `json:"api_key,omitempty"`
	APIBase           string             `json:"api_base,omitempty"`
	Proxy             string             `json:"proxy,omitempty"`
	Timeout           *int               `json:"timeout,omitempty"`
	StreamTimeout     *int               `json:"stream_timeout,omitempty"`
	DownloadMedia     bool               `json:"download_media,omitempty"`
	Raw               bool               `json:"raw,omitempty"`
}

type ChatCompletionResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []Choice       `json:"choices"`
	Usage             *Usage         `json:"usage,omitempty"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Provider          string         `json:"provider,omitempty"`
	ConversationID    string         `json:"conversation_id,omitempty"`
	Conversation      map[string]any `json:"conversation,omitempty"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	FinishReason string    `json:"finish_reason,omitempty"`
	LogProbs     *LogProbs `json:"logprobs,omitempty"`
}

type LogProbs struct {
	Content []TokenLogProb `json:"content,omitempty"`
}

type TokenLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
}

type Usage struct {
	PromptTokens            int           `json:"prompt_tokens"`
	CompletionTokens        int           `json:"completion_tokens"`
	TotalTokens             int           `json:"total_tokens"`
	CompletionTokensDetails *UsageDetails `json:"completion_tokens_details,omitempty"`
}

type UsageDetails struct {
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
	AudioTokens     int `json:"audio_tokens,omitempty"`
	AcceptedTokens  int `json:"accepted_prediction_tokens,omitempty"`
	RejectedTokens  int `json:"rejected_prediction_tokens,omitempty"`
}

type ChatCompletionChunk struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []ChunkChoice  `json:"choices"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Provider          string         `json:"provider,omitempty"`
	ConversationID    string         `json:"conversation_id,omitempty"`
	Conversation      map[string]any `json:"conversation,omitempty"`
}

type ChunkChoice struct {
	Index        int          `json:"index"`
	Delta        MessageDelta `json:"delta"`
	FinishReason string       `json:"finish_reason,omitempty"`
	LogProbs     *LogProbs    `json:"logprobs,omitempty"`
}

type MessageDelta struct {
	Role             string     `json:"role,omitempty"`
	Content          string     `json:"content,omitempty"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
}

type ImageGenerationRequest struct {
	Prompt            string         `json:"prompt"`
	Model             string         `json:"model,omitempty"`
	N                 *int           `json:"n,omitempty"`
	Size              string         `json:"size,omitempty"`
	Quality           string         `json:"quality,omitempty"`
	ResponseFormat    string         `json:"response_format,omitempty"`
	Style             string         `json:"style,omitempty"`
	User              string         `json:"user,omitempty"`
	Messages          []Message      `json:"messages,omitempty"`
	Width             *int           `json:"width,omitempty"`
	Height            *int           `json:"height,omitempty"`
	AspectRatio       string         `json:"aspect_ratio,omitempty"`
	Resolution        string         `json:"resolution,omitempty"`
	NumInferenceSteps *int           `json:"num_inference_steps,omitempty"`
	GuidanceScale     *int           `json:"guidance_scale,omitempty"`
	Seed              *int           `json:"seed,omitempty"`
	NegativePrompt    string         `json:"negative_prompt,omitempty"`
	ImageURL          string         `json:"image_url,omitempty"`
	ImageData         string         `json:"image_data,omitempty"`
	Provider          string         `json:"provider,omitempty"`
	Account           string         `json:"account,omitempty"`
	ConversationID    string         `json:"conversation_id,omitempty"`
	ParentMessageID   string         `json:"parent_message_id,omitempty"`
	Conversation      map[string]any `json:"conversation,omitempty"`
	APIKey            string         `json:"api_key,omitempty"`
	Proxy             string         `json:"proxy,omitempty"`
	Audio             map[string]any `json:"audio,omitempty"`
	DownloadMedia     bool           `json:"download_media,omitempty"`
}

type ImageGenerationResponse struct {
	Created      int64          `json:"created"`
	Data         []ImageData    `json:"data"`
	Provider     string         `json:"provider,omitempty"`
	Conversation map[string]any `json:"conversation,omitempty"`
}

type ImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	Created           int64             `json:"created"`
	OwnedBy           string            `json:"owned_by"`
	Type              string            `json:"type,omitempty"`
	Capabilities      []string          `json:"capabilities,omitempty"`
	Provider          string            `json:"provider,omitempty"`
	Parameters        *ModelParameters  `json:"parameters,omitempty"`
	Category          string            `json:"category,omitempty"`
	Label             string            `json:"label,omitempty"`
	Description       string            `json:"description,omitempty"`
	SubscriptionLevel string            `json:"subscription_level,omitempty"`
	SupportedFeatures []string          `json:"supported_features,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

type ModelParameters struct {
	Schema []ParameterSchema `json:"schema,omitempty"`
}

type ParameterSchema struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Default     any      `json:"default,omitempty"`
	Options     []string `json:"options,omitempty"`
	Minimum     *float64 `json:"minimum,omitempty"`
	Maximum     *float64 `json:"maximum,omitempty"`
	MinLength   *int     `json:"min_length,omitempty"`
	MaxLength   *int     `json:"max_length,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Membership string    `json:"membership"`
	APIKey     string    `json:"api_key,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProviderCredential struct {
	ID                string         `json:"id"`
	Provider          string         `json:"provider"`
	AccountName       string         `json:"account_name"`
	DisplayName       string         `json:"display_name"`
	AccountEmail      string         `json:"account_email"`
	ProviderAccountID string         `json:"provider_account_id"`
	AccountLabel      string         `json:"account_label"`
	Plan              string         `json:"plan"`
	Membership        string         `json:"membership"`
	ProxyURL          string         `json:"proxy_url"`
	ProxyScope        string         `json:"proxy_scope"`
	Enabled           bool           `json:"enabled"`
	LastAuthAt        *time.Time     `json:"last_auth_at,omitempty"`
	CapturedAt        *time.Time     `json:"captured_at,omitempty"`
	Health            *AccountHealth `json:"health,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type AccountHealth struct {
	ID                    string               `json:"id"`
	Provider              string               `json:"provider"`
	AccountName           string               `json:"account_name"`
	AccountEmail          string               `json:"account_email,omitempty"`
	Status                string               `json:"status"`
	CooldownUntil         *time.Time           `json:"cooldown_until,omitempty"`
	ConsecutiveFailures   int                  `json:"consecutive_failures"`
	ConsecutiveRateLimits int                  `json:"consecutive_rate_limits"`
	TotalRequests         int                  `json:"total_requests"`
	TotalSuccesses        int                  `json:"total_successes"`
	TotalFailures         int                  `json:"total_failures"`
	QuotaRemaining        *int                 `json:"quota_remaining,omitempty"`
	QuotaResetAt          *time.Time           `json:"quota_reset_at,omitempty"`
	QuotaFeature          string               `json:"quota_feature,omitempty"`
	LastError             string               `json:"last_error,omitempty"`
	LastRequestAt         *time.Time           `json:"last_request_at,omitempty"`
	LastSuccessAt         *time.Time           `json:"last_success_at,omitempty"`
	LastFailureAt         *time.Time           `json:"last_failure_at,omitempty"`
	Metadata              json.RawMessage      `json:"metadata,omitempty"`
	Stats5m               ProviderRequestStats `json:"stats_5m,omitempty"`
	Stats1h               ProviderRequestStats `json:"stats_1h,omitempty"`
	Stats24h              ProviderRequestStats `json:"stats_24h,omitempty"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

type ProviderRequestLog struct {
	ID           string          `json:"id"`
	Provider     string          `json:"provider"`
	AccountName  string          `json:"account_name"`
	AccountEmail string          `json:"account_email,omitempty"`
	Model        string          `json:"model"`
	RequestType  string          `json:"request_type"`
	Status       string          `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
	DurationMS   int64           `json:"duration_ms"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ProviderRequestStats struct {
	Provider      string    `json:"provider"`
	AccountName   string    `json:"account_name"`
	AccountEmail  string    `json:"account_email,omitempty"`
	DisplayName   string    `json:"display_name,omitempty"`
	Plan          string    `json:"plan,omitempty"`
	Membership    string    `json:"membership,omitempty"`
	RequestType   string    `json:"request_type"`
	Total         int       `json:"total"`
	Success       int       `json:"success"`
	Errors        int       `json:"errors"`
	RateLimited   int       `json:"rate_limited"`
	AuthExpired   int       `json:"auth_expired"`
	Canceled      int       `json:"canceled"`
	AvgDurationMS int64     `json:"avg_duration_ms"`
	LastRequestAt time.Time `json:"last_request_at"`
}
