package rag

import (
	"github.com/tmc/langchaingo/embeddings"
)

type Config struct {
	Embedder  embeddings.Embedder
	Scheme    string
	Host      string
	IndexName string
}

type Option func(*Config)

func WithEmbedder(embedder embeddings.Embedder) Option {
	return func(config *Config) {
		config.Embedder = embedder
	}
}

func WithScheme(scheme string) Option {
	return func(config *Config) {
		config.Scheme = scheme
	}
}

func WithHost(host string) Option {
	return func(config *Config) {
		config.Host = host
	}
}

func WithIndexName(indexName string) Option {
	return func(config *Config) {
		config.IndexName = indexName
	}
}
