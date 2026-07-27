package aigchq

import "context"

const (
	// Current Gemini Web text models.
	ModelGemini35FlashLite = "gemini-3.5-flash-lite"
	ModelGemini36Flash     = "gemini-3.6-flash"
	ModelGemini31Pro       = "gemini-3.1-pro"

	// Explicit extended-thinking model IDs remain available for compatibility.
	// New code can select a base model and set ReasoningEffort instead.
	ModelGemini35FlashLiteExtended = "gemini-3.5-flash-lite-extended"
	ModelGemini36FlashExtended     = "gemini-3.6-flash-extended"
	ModelGemini31ProExtended       = "gemini-3.1-pro-extended"

	ReasoningEffortNone     = "none"
	ReasoningEffortMinimal  = "minimal"
	ReasoningEffortLow      = "low"
	ReasoningEffortMedium   = "medium"
	ReasoningEffortHigh     = "high"
	ReasoningEffortXHigh    = "xhigh"
	ReasoningEffortMax      = "max"
	ReasoningEffortExtended = "extended"
)

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
