package repositories

type OllamaRepository interface {
	Generate(model string, prompt string, stream bool) (string, error)
}
