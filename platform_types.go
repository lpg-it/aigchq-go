package aigchq

import (
	"encoding/json"
	"time"
)

const (
	ProviderChatGPT = "chatgpt-web"
	ProviderGemini  = "gemini-web"
	ProviderQwen    = "qwen-web"
	ImageHostFImage = "f-image"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type userResponse struct {
	User *User `json:"user"`
}

type okResponse struct {
	OK bool `json:"ok"`
}

type ProviderStatus struct {
	Provider     string   `json:"provider"`
	Label        string   `json:"label"`
	Configured   bool     `json:"configured"`
	Accounts     int      `json:"accounts"`
	Capabilities []string `json:"capabilities"`
}

type ProviderCredentialRequest struct {
	Provider          string          `json:"provider"`
	AccountName       string          `json:"account_name,omitempty"`
	DisplayName       string          `json:"display_name,omitempty"`
	AccountEmail      string          `json:"account_email,omitempty"`
	ProviderAccountID string          `json:"provider_account_id,omitempty"`
	AccountLabel      string          `json:"account_label,omitempty"`
	Plan              string          `json:"plan,omitempty"`
	Membership        string          `json:"membership,omitempty"`
	ProxyURL          string          `json:"proxy_url,omitempty"`
	ProxyScope        string          `json:"proxy_scope,omitempty"`
	AuthPayload       json.RawMessage `json:"auth_payload,omitempty"`
	Enabled           *bool           `json:"enabled,omitempty"`
}

type ProviderCredentialSettingsRequest struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	ProxyURL   string `json:"proxy_url,omitempty"`
	ProxyScope string `json:"proxy_scope,omitempty"`
}

type ProviderCredentialMutationResponse struct {
	Item  ProviderCredential   `json:"item"`
	Items []ProviderCredential `json:"items"`
}

type ProviderCredentialsResponse struct {
	Items []ProviderCredential `json:"items"`
}

type ProviderStatusResponse struct {
	Items []ProviderStatus `json:"items"`
}

type ProviderRequestLogsResponse struct {
	Items []ProviderRequestLog `json:"items"`
}

type ProviderRequestStatsResponse struct {
	Items []ProviderRequestStats `json:"items"`
	Hours int                    `json:"hours"`
}

type ImageHostConfig struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Name      string    `json:"name"`
	BaseURL   string    `json:"base_url"`
	Enabled   bool      `json:"enabled"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageHostRequest struct {
	Provider  string `json:"provider"`
	Name      string `json:"name"`
	BaseURL   string `json:"base_url,omitempty"`
	APIToken  string `json:"api_token,omitempty"`
	Enabled   *bool  `json:"enabled,omitempty"`
	IsDefault *bool  `json:"is_default,omitempty"`
}

type ImageHostsResponse struct {
	Items []ImageHostConfig `json:"items"`
}

type ImageHostMutationResponse struct {
	Item  ImageHostConfig   `json:"item"`
	Items []ImageHostConfig `json:"items"`
}

type Conversation struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Model        string          `json:"model"`
	MessageCount int             `json:"message_count"`
	Messages     json.RawMessage `json:"messages,omitempty"`
	PinnedAt     *time.Time      `json:"pinned_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ConversationMessage struct {
	ID             string          `json:"id"`
	ConversationID string          `json:"conversation_id"`
	Role           string          `json:"role"`
	Content        json.RawMessage `json:"content"`
	ContentText    string          `json:"content_text"`
	Position       int             `json:"position"`
	Model          string          `json:"model"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ConversationListOptions struct {
	Limit  int
	Offset int
	Pinned *bool
}

type ConversationsResponse struct {
	Items   []Conversation `json:"items"`
	HasMore bool           `json:"has_more"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
}

type ConversationResponse struct {
	Item Conversation `json:"item"`
}

type ConversationMessagesResponse struct {
	Items []ConversationMessage `json:"items"`
}

type PinConversationRequest struct {
	Pinned bool `json:"pinned"`
}

type AuthCaptureRequest struct {
	Provider    string `json:"provider"`
	AccountName string `json:"account_name,omitempty"`
}

type AuthCapture struct {
	ID          string          `json:"id"`
	Provider    string          `json:"provider"`
	AccountName string          `json:"account_name"`
	Status      string          `json:"status"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	ExpiresAt   time.Time       `json:"expires_at"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

type AuthCaptureResponse struct {
	Capture          AuthCapture `json:"capture"`
	AuthorizationURL string      `json:"authorization_url"`
	CaptureURL       string      `json:"capture_url"`
	CompleteURL      string      `json:"complete_url"`
	ProviderURL      string      `json:"provider_url"`
}

type AuthCaptureGetResponse struct {
	Capture AuthCapture `json:"capture"`
}

type AuthCaptureCompleteResponse struct {
	OK          bool   `json:"ok"`
	Provider    string `json:"provider"`
	AccountName string `json:"account_name"`
}
