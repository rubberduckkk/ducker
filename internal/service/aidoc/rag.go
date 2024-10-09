package aidoc

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"

	"github.com/rubberduckkk/ducker/internal/domain/rag"
	"github.com/rubberduckkk/ducker/internal/domain/rag/valueobj"
)

type Service interface {
	AddDocuments(ctx context.Context, texts []string) error
	QueryDocuments(ctx context.Context, content string, opts ...QueryOption) (*valueobj.QueryResult, error)
}

type svcImpl struct {
	ragRepo  rag.Repository
	llmModel llms.Model
}

func New(ragRepo rag.Repository, llm llms.Model) Service {
	return &svcImpl{ragRepo: ragRepo, llmModel: llm}
}

func (s *svcImpl) AddDocuments(ctx context.Context, texts []string) error {
	docs := make([]schema.Document, 0, len(texts))
	for _, text := range texts {
		docs = append(docs, schema.Document{PageContent: text})
	}
	_, err := s.ragRepo.AddDocuments(ctx, docs)
	return err
}

func (s *svcImpl) QueryDocuments(ctx context.Context, content string, opts ...QueryOption) (*valueobj.QueryResult, error) {
	param := defaultQueryParam()
	for _, fn := range opts {
		fn(param)
	}
	logrus.WithField("content", content).Infof("query docs")
	docs, err := s.ragRepo.SimilaritySearch(ctx, content, param.NumDocuments)
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return &valueobj.QueryResult{
			Summary: "no related documents found",
		}, nil
	}
	logrus.WithField("docs", docs).Infof("similarity search")
	docsContents := make([]string, 0, len(docs))
	for _, doc := range docs {
		docsContents = append(docsContents, doc.PageContent)
	}
	ragQuery := fmt.Sprintf(ragTemplateStr, content, strings.Join(docsContents, "\n"))
	result, err := llms.GenerateFromSinglePrompt(ctx, s.llmModel, ragQuery)
	if err != nil {
		return nil, err
	}
	return &valueobj.QueryResult{
		Summary:     result,
		OriginalDoc: docs[0].PageContent,
	}, nil
}

const ragTemplateStr = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.

Question:
%s

Context:
%s
`
