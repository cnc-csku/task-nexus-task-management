package llm

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/llm"
	"github.com/google/generative-ai-go/genai"
)

type GeminiRepositoryImpl struct {
	wrapperClient *llm.GeminiClient
	cfg           *config.Config
}

func NewGeminiRepo(
	client *llm.GeminiClient,
	cfg *config.Config,
) repositories.GeminiRepository {
	return &GeminiRepositoryImpl{
		wrapperClient: client,
		cfg:           cfg,
	}
}

func (g *GeminiRepositoryImpl) GenerateTaskDescription(ctx context.Context, prompt string) ([]*genai.Content, error) {
	generatedContent, err := g.wrapperClient.Client.GenerateContent(
		ctx,
		genai.Text(
			""+prompt,
		),
	)
	if err != nil {
		return nil, err
	}

	var response []*genai.Content
	for _, content := range generatedContent.Candidates {
		response = append(response, content.Content)
	}

	return response, nil
}
