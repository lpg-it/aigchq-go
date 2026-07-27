package aigchq

import "context"

// ListQwenModels returns models available for the configured Qwen Web account pool.
func (c *Client) ListQwenModels(ctx context.Context) (*ModelsResponse, error) {
	return c.Provider(ProviderQwen).ListModels(ctx)
}

// CreateQwenChatCompletion sends a chat request through the Qwen Web provider route.
func (c *Client) CreateQwenChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return c.Provider(ProviderQwen).CreateChatCompletion(ctx, req)
}

// CreateQwenChatCompletionStream streams a chat response through the Qwen Web provider route.
func (c *Client) CreateQwenChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	return c.Provider(ProviderQwen).CreateChatCompletionStream(ctx, req)
}

// CreateAsyncQwenChatCompletion creates an async Qwen Web chat task.
func (c *Client) CreateAsyncQwenChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*AsyncTaskResponse, error) {
	return c.Provider(ProviderQwen).CreateAsyncChatCompletion(ctx, req)
}

// CreateQwenChatCompletionAndWait creates an async Qwen Web chat task and waits for completion.
func (c *Client) CreateQwenChatCompletionAndWait(ctx context.Context, req *ChatCompletionRequest, opts ...WaitOption) (*ChatCompletionResponse, error) {
	return c.Provider(ProviderQwen).CreateChatCompletionAndWait(ctx, req, opts...)
}
