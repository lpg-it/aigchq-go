package aigchq

import "context"

func (c *Client) ListModels(ctx context.Context) (*ModelsResponse, error) {
	var out ModelsResponse
	if err := c.doJSON(ctx, httpMethodGet, "/v1/models", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetModel(ctx context.Context, model string) (*Model, error) {
	var out Model
	if err := c.doJSON(ctx, httpMethodGet, "/v1/models/"+pathEscape(model), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var out HealthResponse
	if err := c.doJSON(ctx, httpMethodGet, "/health", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
