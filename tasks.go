package aigchq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type TaskType string

const (
	TaskTypeChat  TaskType = "chat"
	TaskTypeImage TaskType = "image"
	TaskTypeVideo TaskType = "video"
)

type AsyncTaskResponse struct {
	TaskID    string     `json:"task_id"`
	Status    TaskStatus `json:"status"`
	PollURL   string     `json:"poll_url,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type TaskResponse struct {
	TaskID      string          `json:"task_id"`
	Provider    string          `json:"provider"`
	Model       string          `json:"model"`
	Type        TaskType        `json:"type"`
	Status      TaskStatus      `json:"status"`
	Progress    float64         `json:"progress,omitempty"`
	Result      json.RawMessage `json:"result,omitempty"`
	Error       string          `json:"error,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

func (t *TaskResponse) IsTerminal() bool {
	return t != nil && (t.Status == TaskStatusCompleted || t.Status == TaskStatusFailed)
}

func (t *TaskResponse) ChatResult() (*ChatCompletionResponse, error) {
	if t == nil {
		return nil, errors.New("task is nil")
	}
	if t.Status != TaskStatusCompleted {
		return nil, fmt.Errorf("task is not completed: %s", t.Status)
	}
	var out ChatCompletionResponse
	if err := json.Unmarshal(t.Result, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (t *TaskResponse) ImageResult() (*ImageGenerationResponse, error) {
	if t == nil {
		return nil, errors.New("task is nil")
	}
	if t.Status != TaskStatusCompleted {
		return nil, fmt.Errorf("task is not completed: %s", t.Status)
	}
	var out ImageGenerationResponse
	if err := json.Unmarshal(t.Result, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type WaitOptions struct {
	Interval time.Duration
	Timeout  time.Duration
	OnPoll   func(*TaskResponse)
}

type WaitOption func(*WaitOptions)

func WithPollInterval(interval time.Duration) WaitOption {
	return func(o *WaitOptions) {
		o.Interval = interval
	}
}

func WithPollTimeout(timeout time.Duration) WaitOption {
	return func(o *WaitOptions) {
		o.Timeout = timeout
	}
}

func WithPollHook(fn func(*TaskResponse)) WaitOption {
	return func(o *WaitOptions) {
		o.OnPoll = fn
	}
}

func defaultWaitOptions(opts ...WaitOption) WaitOptions {
	out := WaitOptions{
		Interval: 2 * time.Second,
		Timeout:  30 * time.Minute,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&out)
		}
	}
	if out.Interval <= 0 {
		out.Interval = 2 * time.Second
	}
	if out.Timeout <= 0 {
		out.Timeout = 30 * time.Minute
	}
	return out
}

func (c *Client) GetTask(ctx context.Context, taskID string) (*TaskResponse, error) {
	if taskID == "" {
		return nil, errors.New("task ID is required")
	}
	var out TaskResponse
	if err := c.doJSON(ctx, httpMethodGet, "/v1/tasks/"+pathEscape(taskID), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) WaitTask(ctx context.Context, taskID string, opts ...WaitOption) (*TaskResponse, error) {
	cfg := defaultWaitOptions(opts...)
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	for {
		task, err := c.GetTask(ctx, taskID)
		if err != nil {
			return nil, err
		}
		if cfg.OnPoll != nil {
			cfg.OnPoll(task)
		}
		switch task.Status {
		case TaskStatusCompleted:
			return task, nil
		case TaskStatusFailed:
			if task.Error == "" {
				return task, errors.New("task failed")
			}
			return task, errors.New(task.Error)
		}
		if err := sleepContext(ctx, cfg.Interval); err != nil {
			return task, err
		}
	}
}

func (c *Client) WaitChatCompletion(ctx context.Context, taskID string, opts ...WaitOption) (*ChatCompletionResponse, error) {
	task, err := c.WaitTask(ctx, taskID, opts...)
	if err != nil {
		return nil, err
	}
	return task.ChatResult()
}

func (c *Client) WaitImageGeneration(ctx context.Context, taskID string, opts ...WaitOption) (*ImageGenerationResponse, error) {
	task, err := c.WaitTask(ctx, taskID, opts...)
	if err != nil {
		return nil, err
	}
	return task.ImageResult()
}
