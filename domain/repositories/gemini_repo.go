package repositories

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

type GeminiRepository interface {
	GenerateTaskDescription(ctx context.Context, prompt string) ([]*genai.Content, error)
}
