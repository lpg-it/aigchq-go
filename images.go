package aigchq

import "context"

func (c *Client) CreateImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	if req == nil {
		req = &ImageGenerationRequest{}
	}
	var out ImageGenerationResponse
	if err := c.doJSON(ctx, httpMethodPost, "/v1/images/generations", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateAsyncImageGeneration(ctx context.Context, req *ImageGenerationRequest) (*AsyncTaskResponse, error) {
	if req == nil {
		req = &ImageGenerationRequest{}
	}
	var out AsyncTaskResponse
	if err := c.doJSON(ctx, httpMethodPost, "/v1/async/images/generations", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateImageGenerationAndWait(ctx context.Context, req *ImageGenerationRequest, opts ...WaitOption) (*ImageGenerationResponse, error) {
	task, err := c.CreateAsyncImageGeneration(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.WaitImageGeneration(ctx, task.TaskID, opts...)
}
