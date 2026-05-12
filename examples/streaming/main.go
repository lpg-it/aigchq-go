package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	aigchq "github.com/lpg-it/aigchq-go"
)

func main() {
	client := mustClient()

	stream, err := client.CreateChatCompletionStream(context.Background(), &aigchq.ChatCompletionRequest{
		Model: "gpt-5-5-instant",
		Messages: []aigchq.Message{
			{Role: "user", Content: "用三点说明同步流式接口适合什么场景。"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			fmt.Println()
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
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
