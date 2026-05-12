package main

import (
	"context"
	"fmt"
	"log"
	"os"

	aigchq "github.com/lpg-it/aigchq-go"
)

func main() {
	client := mustClient()

	resp, err := client.CreateChatCompletion(context.Background(), &aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-thinking",
		Messages: []aigchq.Message{
			{Role: "user", Content: "用一句话介绍 AIGCHQ。"},
		},
	})
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
