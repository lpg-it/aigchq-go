package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aigchq "github.com/lpg-it/aigchq-go"
)

func main() {
	client := mustClient()
	ctx := context.Background()

	task, err := client.CreateAsyncChatCompletion(ctx, &aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-thinking",
		Messages: []aigchq.Message{
			{Role: "user", Content: "给我一个适合异步任务的使用场景。"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("task_id=%s status=%s poll_url=%s\n", task.TaskID, task.Status, task.PollURL)

	resp, err := client.WaitChatCompletion(
		ctx,
		task.TaskID,
		aigchq.WithPollInterval(2*time.Second),
		aigchq.WithPollTimeout(30*time.Minute),
		aigchq.WithPollHook(func(task *aigchq.TaskResponse) {
			fmt.Printf("poll status=%s progress=%.0f%%\n", task.Status, task.Progress*100)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

func mustClient() *aigchq.Client {
	apiKey := os.Getenv("AIGCHQ_API_KEY")
	if apiKey == "" {
		log.Fatal("AIGCHQ_API_KEY is required")
	}
	opts := []aigchq.Option{}
	if baseURL := os.Getenv("AIGCHQ_BASE_URL"); baseURL != "" {
		opts = append(opts, aigchq.WithBaseURL(baseURL))
	}
	client, err := aigchq.NewClient(apiKey, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
