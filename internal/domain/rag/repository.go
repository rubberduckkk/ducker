package rag

import (
	"context"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type Repository interface {
	AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error)
	SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error)
}
