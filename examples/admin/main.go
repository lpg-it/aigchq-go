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

	credentials, err := client.ListProviderCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("provider accounts: %d\n", len(credentials))
	for _, item := range credentials {
		fmt.Printf("%s %s plan=%s enabled=%v\n", item.Provider, item.AccountName, item.Plan, item.Enabled)
	}

	logs, err := client.ListProviderRequestLogs(ctx, 10)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range logs {
		fmt.Printf("%s %s %s %s %s\n", item.CreatedAt.Format(time.RFC3339), item.Provider, item.Model, item.Status, item.ErrorMessage)
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
