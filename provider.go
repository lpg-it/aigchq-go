package aigchq

import "context"

type ProviderClient struct {
	client   *Client
	provider string
}

func (c *Client) Provider(provider string) *ProviderClient {
	return &ProviderClient{client: c, provider: provider}
}

func (p *ProviderClient) ListModels(ctx context.Context) (*ModelsResponse, error) {
	var out ModelsResponse
	if err := p.client.doJSON(ctx, httpMethodGet, "/api/"+pathEscape(p.provider)+"/models", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *ProviderClient) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	var out ChatCompletionResponse
	if err := p.client.doJSON(ctx, httpMethodPost, "/api/"+pathEscape(p.provider)+"/chat/completions", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *ProviderClient) CreateAsyncChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*AsyncTaskResponse, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	copyReq := *req
	copyReq.Stream = false
	var out AsyncTaskResponse
	if err := p.client.doJSON(ctx, httpMethodPost, "/api/"+pathEscape(p.provider)+"/async/chat/completions", nil, &copyReq, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *ProviderClient) CreateChatCompletionAndWait(ctx context.Context, req *ChatCompletionRequest, opts ...WaitOption) (*ChatCompletionResponse, error) {
	task, err := p.CreateAsyncChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	return p.client.WaitChatCompletion(ctx, task.TaskID, opts...)
}

func (p *ProviderClient) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	if req == nil {
		req = &ChatCompletionRequest{}
	}
	copyReq := *req
	copyReq.Stream = true
	return p.client.createChatStream(ctx, "/api/"+pathEscape(p.provider)+"/chat/completions", &copyReq)
}

func (p *ProviderClient) CreateImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	if req == nil {
		req = &ImageGenerationRequest{}
	}
	var out ImageGenerationResponse
	if err := p.client.doJSON(ctx, httpMethodPost, "/api/"+pathEscape(p.provider)+"/images/generations", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *ProviderClient) CreateAsyncImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*AsyncTaskResponse, error) {
	if req == nil {
		req = &ImageGenerationRequest{}
	}
	var out AsyncTaskResponse
	if err := p.client.doJSON(ctx, httpMethodPost, "/api/"+pathEscape(p.provider)+"/async/images/generations", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *ProviderClient) CreateImageGenerationAndWait(ctx context.Context, req *ImageGenerationRequest, opts ...WaitOption) (*ImageGenerationResponse, error) {
	task, err := p.CreateAsyncImageGeneration(ctx, req)
	if err != nil {
		return nil, err
	}
	return p.client.WaitImageGeneration(ctx, task.TaskID, opts...)
}
