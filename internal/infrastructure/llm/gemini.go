package llm

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	Client *genai.GenerativeModel
}

func NewGeminiClient(
	ctx context.Context,
	cfg *config.Config,
) *GeminiClient {
	client, err := genai.NewClient(
		ctx,
		option.WithAPIKey(cfg.GeminiClient.ApiKey),
	)
	if err != nil {
		log.Fatalf("❌ failed to create Gemini client: %v", err)
		return nil
	}

	if cfg.GeminiClient.Model == "" {
		cfg.GeminiClient.Model = "gemini-2.0-flash"
	}

	model := client.GenerativeModel(cfg.GeminiClient.Model)

	return &GeminiClient{Client: model}
}
