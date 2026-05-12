package aigchq

import "context"

func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	var out ChatCompletionResponse
	if err := c.doJSON(ctx, httpMethodPost, "/v1/chat/completions", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateAsyncChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*AsyncTaskResponse, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	copyReq := *req
	copyReq.Stream = false
	var out AsyncTaskResponse
	if err := c.doJSON(ctx, httpMethodPost, "/v1/async/chat/completions", nil, &copyReq, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateChatCompletionAndWait(ctx context.Context, req *ChatCompletionRequest, opts ...WaitOption) (*ChatCompletionResponse, error) {
	task, err := c.CreateAsyncChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.WaitChatCompletion(ctx, task.TaskID, opts...)
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	copyReq := *req
	copyReq.Stream = true
	return c.createChatStream(ctx, "/v1/chat/completions", &copyReq)
}
