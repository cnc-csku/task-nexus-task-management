package llm

import (
	"bytes"
	"encoding/json"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/llm"
)

type OllamaRepositoryImpl struct {
	client *llm.OllamaClient
	cfg    *config.Config
}

func NewOllamaRepository(client *llm.OllamaClient, cfg *config.Config) repositories.OllamaRepository {
	return &OllamaRepositoryImpl{
		client: client,
		cfg:    cfg,
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Response           string        `json:"response"`
	Done               bool          `json:"done"`
	DoneReason         string        `json:"done_reason"`
	Context            []interface{} `json:"context"`
	TotalDuration      int64         `json:"total_duration"`
	LoadDuration       int64         `json:"load_duration"`
	PromptEvalCount    int64         `json:"prompt_eval_count"`
	PromptEvalDuration int64         `json:"prompt_eval_duration"`
	EvalCount          int64         `json:"eval_count"`
	EvalDuration       int64         `json:"eval_duration"`
}

func (r *OllamaRepositoryImpl) Generate(model string, prompt string, stream bool) (string, error) {
	request := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: stream,
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	endpoint := "http://" + r.cfg.OllamaClient.Endpoint + "/api/generate"
	response, err := r.client.HTTPClient.Post(
		endpoint,
		"application/json",
		bytes.NewBuffer(requestJson),
	)
	if err != nil {
		return "", err
	}

	var ollamaResponse OllamaResponse
	err = json.NewDecoder(response.Body).Decode(&ollamaResponse)
	if err != nil {
		return "", err
	}

	return ollamaResponse.Response, nil
}
