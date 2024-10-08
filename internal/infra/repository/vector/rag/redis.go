package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/redisvector"

	"github.com/rubberduckkk/ducker/internal/domain/rag"
)

type redisRepository struct {
	store *redisvector.Store
}

func NewRedisRepo(opts ...Option) (rag.Repository, error) {
	cfg := new(Config)
	for _, opt := range opts {
		opt(cfg)
	}

	url := fmt.Sprintf("%v://%v", cfg.Scheme, cfg.Host)
	store, err := redisvector.New(context.Background(),
		redisvector.WithConnectionURL(url),
		redisvector.WithIndexName(cfg.IndexName, true),
		redisvector.WithEmbedder(cfg.Embedder),
	)
	if err != nil {
		return nil, err
	}
	return &redisRepository{store: store}, nil
}

func (r *redisRepository) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	return r.store.AddDocuments(ctx, docs)
}

func (r *redisRepository) SimilaritySearch(
	ctx context.Context,
	query string,
	numDocuments int,
	options ...vectorstores.Option) ([]schema.Document, error) {
	return r.store.SimilaritySearch(ctx, query, numDocuments, options...)
}
