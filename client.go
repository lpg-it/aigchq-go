package aigchq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://aigchq.com"
)

// Client is the official Go client for the AIGCHQ API.
//
// The standard OpenAI-compatible chat/image endpoints are synchronous.
// Use CreateAsyncChatCompletion/CreateAsyncImageGeneration plus Wait* helpers
// when you want request/response lifetimes decoupled from upstream execution.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
	headers    http.Header
	retry      RetryConfig
}

// Option configures a Client.
type Option func(*Client) error

// RetryConfig controls automatic retry behavior for transient transport, 429,
// and 5xx failures. Retries are disabled when MaxRetries is 0.
type RetryConfig struct {
	MaxRetries int
	MinDelay   time.Duration
	MaxDelay   time.Duration
}

// NewClient creates a client using the provided API key.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	u, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		baseURL: u,
		apiKey:  strings.TrimSpace(apiKey),
		headers: make(http.Header),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		retry: RetryConfig{
			MaxRetries: 2,
			MinDelay:   500 * time.Millisecond,
			MaxDelay:   5 * time.Second,
		},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// MustNewClient is like NewClient but panics on configuration errors.
func MustNewClient(apiKey string, opts ...Option) *Client {
	c, err := NewClient(apiKey, opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// WithBaseURL sets the API base URL. Example: https://aigchq.com.
func WithBaseURL(raw string) Option {
	return func(c *Client) error {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return errors.New("base URL is required")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid base URL: %s", raw)
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient injects a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return errors.New("http client is nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout sets the default timeout on the SDK-owned HTTP client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		if timeout <= 0 {
			return errors.New("timeout must be positive")
		}
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithHeader adds a header to every request.
func WithHeader(key, value string) Option {
	return func(c *Client) error {
		key = strings.TrimSpace(key)
		if key == "" {
			return errors.New("header key is required")
		}
		c.headers.Set(key, value)
		return nil
	}
}

// WithRetry configures retry behavior.
func WithRetry(retry RetryConfig) Option {
	return func(c *Client) error {
		if retry.MaxRetries < 0 {
			return errors.New("max retries cannot be negative")
		}
		if retry.MinDelay <= 0 {
			retry.MinDelay = 500 * time.Millisecond
		}
		if retry.MaxDelay <= 0 {
			retry.MaxDelay = 5 * time.Second
		}
		if retry.MaxDelay < retry.MinDelay {
			return errors.New("max retry delay cannot be smaller than min retry delay")
		}
		c.retry = retry
		return nil
	}
}

// WithNoRetry disables retries.
func WithNoRetry() Option {
	return WithRetry(RetryConfig{})
}

// SetAPIKey updates the API key used for future requests.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = strings.TrimSpace(apiKey)
}

// DoJSON sends a JSON request to an AIGCHQ API path.
//
// This is intended for newly added platform endpoints before the SDK grows a
// first-class method. Normal chat, image, task, provider, and account APIs have
// typed helpers and should use those helpers instead.
func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	return c.doJSON(ctx, method, path, query, body, out)
}

func (c *Client) buildURL(path string, query url.Values) string {
	u := *c.baseURL
	basePath := strings.TrimRight(u.Path, "/")
	reqPath := "/" + strings.TrimLeft(path, "/")
	u.Path = basePath + reqPath
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, []byte, error) {
	var payload []byte
	var reader io.Reader
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(path, query), reader)
	if err != nil {
		return nil, nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("x-api-key", c.apiKey)
	}
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, payload, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	req, payload, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	resp, err := c.do(req, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, out)
}

func (c *Client) do(req *http.Request, payload []byte) (*http.Response, error) {
	attempts := c.retry.MaxRetries + 1
	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if payload != nil {
			req.Body = io.NopCloser(bytes.NewReader(payload))
			req.GetBody = func() (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(payload)), nil
			}
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if !c.shouldRetry(nil, err, attempt) {
				return nil, err
			}
			if err := sleepContext(req.Context(), c.retryDelay(attempt+1, nil)); err != nil {
				return nil, err
			}
			continue
		}
		if !c.shouldRetry(resp, nil, attempt) {
			return resp, nil
		}
		lastErr = apiErrorFromResponse(resp)
		delay := c.retryDelay(attempt+1, resp)
		_ = resp.Body.Close()
		if err := sleepContext(req.Context(), delay); err != nil {
			return nil, err
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("request failed")
}

func (c *Client) shouldRetry(resp *http.Response, err error, attempt int) bool {
	if attempt >= c.retry.MaxRetries {
		return false
	}
	if err != nil {
		return true
	}
	if resp == nil {
		return false
	}
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
}

func (c *Client) retryDelay(attempt int, resp *http.Response) time.Duration {
	if resp != nil {
		if value := strings.TrimSpace(resp.Header.Get("Retry-After")); value != "" {
			if secs, err := strconv.Atoi(value); err == nil && secs >= 0 {
				return time.Duration(secs) * time.Second
			}
			if when, err := http.ParseTime(value); err == nil {
				if delay := time.Until(when); delay > 0 {
					return delay
				}
			}
		}
	}
	minDelay := c.retry.MinDelay
	if minDelay <= 0 {
		minDelay = 500 * time.Millisecond
	}
	maxDelay := c.retry.MaxDelay
	if maxDelay <= 0 {
		maxDelay = 5 * time.Second
	}
	delay := minDelay << max(0, attempt-1)
	if delay > maxDelay {
		delay = maxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(delay / 2)))
	return delay/2 + jitter
}

func decodeResponse(resp *http.Response, out any) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return apiErrorFromResponse(resp)
	}
	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		return err
	}
	return nil
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
