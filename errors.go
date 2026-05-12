package aigchq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// APIError is returned for non-2xx API responses.
type APIError struct {
	StatusCode int
	Status     string
	RequestID  string
	Detail     *ErrorDetail
	Body       []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Detail != nil && strings.TrimSpace(e.Detail.Message) != "" {
		if e.Detail.Code != "" {
			return fmt.Sprintf("aigchq: %s: %s (%s)", e.Status, e.Detail.Message, e.Detail.Code)
		}
		return fmt.Sprintf("aigchq: %s: %s", e.Status, e.Detail.Message)
	}
	if len(e.Body) > 0 {
		return fmt.Sprintf("aigchq: %s: %s", e.Status, strings.TrimSpace(string(e.Body)))
	}
	return "aigchq: " + e.Status
}

func apiErrorFromResponse(resp *http.Response) error {
	if resp == nil {
		return &APIError{StatusCode: 0, Status: "no response"}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		RequestID:  resp.Header.Get("X-Request-Id"),
		Body:       body,
	}
	var wrapped ErrorResponse
	if len(body) > 0 && json.Unmarshal(body, &wrapped) == nil {
		apiErr.Detail = &wrapped.Error
	}
	return apiErr
}
