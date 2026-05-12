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

	resp, err := client.CreateImageGenerationAndWait(
		context.Background(),
		&aigchq.ImageGenerationRequest{
			Model:  "gemini-3.1-pro",
			Prompt: "一张干净的产品发布会海报，中心是 AIGCHQ API 控制台界面。",
			Size:   "1024x1024",
		},
		aigchq.WithPollInterval(3*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range resp.Data {
		fmt.Println(image.URL)
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
