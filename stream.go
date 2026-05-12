package aigchq

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// ChatCompletionStream reads OpenAI-compatible server-sent chat chunks.
type ChatCompletionStream struct {
	body   io.ReadCloser
	reader *bufio.Reader
	done   bool
}

func (c *Client) createChatStream(ctx context.Context, path string, req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	httpReq, payload, err := c.newRequest(ctx, httpMethodPost, path, nil, req)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.do(httpReq, payload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiErrorFromResponse(resp)
	}
	return &ChatCompletionStream{
		body:   resp.Body,
		reader: bufio.NewReader(resp.Body),
	}, nil
}

// Recv returns the next stream chunk. It returns io.EOF after the [DONE] event.
func (s *ChatCompletionStream) Recv() (*ChatCompletionChunk, error) {
	if s == nil {
		return nil, errors.New("stream is nil")
	}
	if s.done {
		return nil, io.EOF
	}
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			s.done = true
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if err != nil {
				s.done = true
				return nil, err
			}
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			if err != nil {
				s.done = true
				return nil, err
			}
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			s.done = true
			return nil, io.EOF
		}
		var chunk ChatCompletionChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return nil, err
		}
		return &chunk, nil
	}
}

// Close closes the underlying HTTP response body.
func (s *ChatCompletionStream) Close() error {
	if s == nil || s.body == nil {
		return nil
	}
	s.done = true
	return s.body.Close()
}
