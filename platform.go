package aigchq

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

func (c *Client) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	if req == nil {
		return nil, errors.New("register request is required")
	}
	var out userResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/auth/register", nil, req, &out); err != nil {
		return nil, err
	}
	if out.User != nil && out.User.APIKey != "" {
		c.SetAPIKey(out.User.APIKey)
	}
	return out.User, nil
}

func (c *Client) Login(ctx context.Context, req *LoginRequest) (*User, error) {
	if req == nil {
		return nil, errors.New("login request is required")
	}
	var out userResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/auth/login", nil, req, &out); err != nil {
		return nil, err
	}
	return out.User, nil
}

func (c *Client) Logout(ctx context.Context) error {
	var out okResponse
	return c.doJSON(ctx, httpMethodPost, "/api/auth/logout", nil, map[string]any{}, &out)
}

func (c *Client) Me(ctx context.Context) (*User, error) {
	var out userResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/auth/me", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.User, nil
}

func (c *Client) UpdateMe(ctx context.Context, req *UpdateUserRequest) (*User, error) {
	if req == nil {
		return nil, errors.New("update user request is required")
	}
	var out userResponse
	if err := c.doJSON(ctx, httpMethodPatch, "/api/auth/me", nil, req, &out); err != nil {
		return nil, err
	}
	return out.User, nil
}

func (c *Client) RegenerateAPIKey(ctx context.Context) (*User, error) {
	var out userResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/auth/api-key", nil, map[string]any{}, &out); err != nil {
		return nil, err
	}
	if out.User != nil && out.User.APIKey != "" {
		c.SetAPIKey(out.User.APIKey)
	}
	return out.User, nil
}

func (c *Client) ProviderStatus(ctx context.Context) ([]ProviderStatus, error) {
	var out ProviderStatusResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/provider-status", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) ListProviderCredentials(ctx context.Context) ([]ProviderCredential, error) {
	var out ProviderCredentialsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/provider-credentials", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) UpsertProviderCredential(ctx context.Context, req *ProviderCredentialRequest) (*ProviderCredentialMutationResponse, error) {
	if req == nil {
		return nil, errors.New("provider credential request is required")
	}
	var out ProviderCredentialMutationResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/provider-credentials", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UpdateProviderCredentialSettings(ctx context.Context, id string, req *ProviderCredentialSettingsRequest) (*ProviderCredentialMutationResponse, error) {
	if id == "" {
		return nil, errors.New("provider credential id is required")
	}
	if req == nil {
		return nil, errors.New("provider credential settings request is required")
	}
	var out ProviderCredentialMutationResponse
	path := "/api/provider-credentials/" + pathEscape(id) + "/settings"
	if err := c.doJSON(ctx, httpMethodPatch, path, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteProviderCredential(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("provider credential id is required")
	}
	var out okResponse
	return c.doJSON(ctx, httpMethodDelete, "/api/provider-credentials/"+pathEscape(id), nil, nil, &out)
}

func (c *Client) ListProviderRequestLogs(ctx context.Context, limit int) ([]ProviderRequestLog, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	var out ProviderRequestLogsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/provider-request-logs", query, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) ListProviderRequestStats(ctx context.Context, hours int) (*ProviderRequestStatsResponse, error) {
	query := url.Values{}
	if hours > 0 {
		query.Set("hours", strconv.Itoa(hours))
	}
	var out ProviderRequestStatsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/provider-request-stats", query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListImageHosts(ctx context.Context) ([]ImageHostConfig, error) {
	var out ImageHostsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/image-hosts", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) UpsertImageHost(ctx context.Context, req *ImageHostRequest) (*ImageHostMutationResponse, error) {
	if req == nil {
		return nil, errors.New("image host request is required")
	}
	var out ImageHostMutationResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/image-hosts", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteImageHost(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("image host id is required")
	}
	var out okResponse
	return c.doJSON(ctx, httpMethodDelete, "/api/image-hosts/"+pathEscape(id), nil, nil, &out)
}

func (c *Client) ListConversations(ctx context.Context, opts ConversationListOptions) (*ConversationsResponse, error) {
	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Pinned != nil {
		query.Set("pinned", strconv.FormatBool(*opts.Pinned))
	}
	var out ConversationsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/conversations", query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListConversationMessages(ctx context.Context, conversationID string) ([]ConversationMessage, error) {
	if conversationID == "" {
		return nil, errors.New("conversation id is required")
	}
	var out ConversationMessagesResponse
	path := "/api/conversations/" + pathEscape(conversationID) + "/messages"
	if err := c.doJSON(ctx, httpMethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) UpsertConversation(ctx context.Context, conv *Conversation) (*Conversation, error) {
	if conv == nil {
		return nil, errors.New("conversation is required")
	}
	var out ConversationResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/conversations", nil, conv, &out); err != nil {
		return nil, err
	}
	return &out.Item, nil
}

func (c *Client) PinConversation(ctx context.Context, conversationID string, pinned bool) (*Conversation, error) {
	if conversationID == "" {
		return nil, errors.New("conversation id is required")
	}
	var out ConversationResponse
	path := "/api/conversations/" + pathEscape(conversationID) + "/pin"
	if err := c.doJSON(ctx, httpMethodPatch, path, nil, PinConversationRequest{Pinned: pinned}, &out); err != nil {
		return nil, err
	}
	return &out.Item, nil
}

func (c *Client) DeleteConversation(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("conversation id is required")
	}
	var out okResponse
	return c.doJSON(ctx, httpMethodDelete, "/api/conversations/"+pathEscape(conversationID), nil, nil, &out)
}

func (c *Client) CreateAuthCapture(ctx context.Context, req *AuthCaptureRequest) (*AuthCaptureResponse, error) {
	if req == nil {
		return nil, errors.New("auth capture request is required")
	}
	var out AuthCaptureResponse
	if err := c.doJSON(ctx, httpMethodPost, "/api/auth-captures", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetAuthCapture(ctx context.Context, id string) (*AuthCapture, error) {
	if id == "" {
		return nil, errors.New("auth capture id is required")
	}
	var out AuthCaptureGetResponse
	if err := c.doJSON(ctx, httpMethodGet, "/api/auth-captures/"+pathEscape(id), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Capture, nil
}

func (c *Client) CompleteAuthCapture(ctx context.Context, id string, payload any) (*AuthCaptureCompleteResponse, error) {
	if id == "" {
		return nil, errors.New("auth capture id is required")
	}
	if payload == nil {
		payload = json.RawMessage(`{}`)
	}
	var out AuthCaptureCompleteResponse
	path := "/api/auth-captures/" + pathEscape(id) + "/complete"
	if err := c.doJSON(ctx, httpMethodPost, path, nil, payload, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
