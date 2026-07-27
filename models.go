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

	// Qwen Web fallback chat models. Live account catalogs from
	// GET /api/qwen-web/models are authoritative when available.
	ModelQwen37Plus  = "qwen3.7-plus"
	ModelQwen36Plus  = "qwen3.6-plus"
	ModelQwen35Flash = "qwen3.5-flash"

	// Explicit Qwen mode suffixes. Prefer a base model plus ReasoningEffort /
	// WebSearch when possible; these aliases remain convenient for clients.
	ModelQwen37PlusThinking       = "qwen3.7-plus-thinking"
	ModelQwen37PlusSearch         = "qwen3.7-plus-search"
	ModelQwen37PlusThinkingSearch = "qwen3.7-plus-thinking-search"
	ModelQwen37PlusFast           = "qwen3.7-plus-fast"
	ModelQwen37PlusFastSearch     = "qwen3.7-plus-fast-search"

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
