package rag

import (
	"context"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/weaviate"

	"github.com/rubberduckkk/ducker/internal/domain/rag"
)

type weaviateRepository struct {
	store *weaviate.Store
}

func NewWeaviateRepo(opts ...Option) (rag.Repository, error) {
	cfg := new(Config)
	for _, opt := range opts {
		opt(cfg)
	}
	store, err := weaviate.New(
		weaviate.WithEmbedder(cfg.Embedder),
		weaviate.WithScheme(cfg.Scheme),
		weaviate.WithHost(cfg.Host),
		weaviate.WithIndexName(cfg.IndexName),
	)
	if err != nil {
		return nil, err
	}
	return &weaviateRepository{store: &store}, nil
}

func (w *weaviateRepository) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	return w.store.AddDocuments(ctx, docs)
}

func (w *weaviateRepository) SimilaritySearch(
	ctx context.Context,
	query string,
	numDocuments int,
	options ...vectorstores.Option) ([]schema.Document, error) {
	return w.store.SimilaritySearch(ctx, query, numDocuments, options...)
}
