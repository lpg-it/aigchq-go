package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	aigchq "github.com/lpg-it/aigchq-go"
)

func main() {
	apiKey := os.Getenv("AIGCHQ_API_KEY")
	if apiKey == "" {
		log.Fatal("AIGCHQ_API_KEY is required")
	}

	opts := []aigchq.Option{aigchq.WithTimeout(20 * time.Minute)}
	if baseURL := os.Getenv("AIGCHQ_BASE_URL"); baseURL != "" {
		opts = append(opts, aigchq.WithBaseURL(baseURL))
	}
	client, err := aigchq.NewClient(apiKey, opts...)
	if err != nil {
		log.Fatal(err)
	}

	model := os.Getenv("AIGCHQ_QWEN_MODEL")
	if model == "" {
		model = aigchq.ModelQwen37PlusThinking
	}

	req := &aigchq.ChatCompletionRequest{
		Model: model,
		Messages: []aigchq.Message{{
			Role:    "user",
			Content: "用一句话介绍 Qwen Web 通过 AIGCHQ 的能力。",
		}},
	}

	if videoPath := os.Getenv("AIGCHQ_QWEN_VIDEO_PATH"); videoPath != "" {
		video, readErr := os.ReadFile(videoPath)
		if readErr != nil {
			log.Fatalf("read video: %v", readErr)
		}
		req.Messages = []aigchq.Message{{
			Role: "user",
			Content: []aigchq.ContentPart{
				{
					Type: "text",
					Text: "请分析附件中的视频，用中文给出可见场景和字幕摘要。",
				},
				{
					Type: "input_file",
					File: &aigchq.InputFile{
						Data:     "data:video/mp4;base64," + base64.StdEncoding.EncodeToString(video),
						Filename: filepath.Base(videoPath),
						MimeType: "video/mp4",
					},
				},
			},
		}}
	}

	resp, err := client.CreateQwenChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Choices) == 0 {
		log.Fatal("empty response")
	}
	if resp.Choices[0].Message.ReasoningContent != "" {
		fmt.Println("thinking:", resp.Choices[0].Message.ReasoningContent)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
