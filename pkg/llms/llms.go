package llms

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/rubberduckkk/ducker/internal/infra/config"
)

type LLM struct {
	Embedder embeddings.Embedder
	OpenAI   *openai.LLM
}

var (
	llmInstance *LLM
)

func Init(cfg config.LLM) error {
	openAI, err := openai.New(
		openai.WithToken(cfg.OpenAIAPIKey),
		openai.WithModel(cfg.OpenAIModelName),
		openai.WithEmbeddingModel(cfg.EmbeddingModelName))
	if err != nil {
		return err
	}

	emb, err := embeddings.NewEmbedder(openAI)
	if err != nil {
		return err
	}

	llmInstance = &LLM{
		Embedder: emb,
		OpenAI:   openAI,
	}
	return nil
}

func Instance() *LLM {
	return llmInstance
}
