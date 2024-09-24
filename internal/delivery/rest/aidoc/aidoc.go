package aidoc

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/rubberduckkk/ducker/internal/delivery/rest"
	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/internal/infra/repository/vector/rag"
	"github.com/rubberduckkk/ducker/internal/service/aidoc"
	"github.com/rubberduckkk/ducker/pkg/llms"
)

type Delivery struct {
	svc aidoc.Service
}

var (
	r     *Delivery
	rOnce sync.Once
)

func Deliver() *Delivery {
	rOnce.Do(func() {
		weav := config.Get().LLM.Weaviate
		ragRepo, err := rag.NewWeaviateRepo(
			rag.WithEmbedder(llms.Instance().Embedder),
			rag.WithScheme(weav.Scheme),
			rag.WithHost(weav.Host),
			rag.WithIndexName(weav.Index),
		)
		if err != nil {
			panic(err)
		}
		svc := aidoc.New(ragRepo, llms.Instance().OpenAI)
		r = &Delivery{svc: svc}
	})
	return r
}

func (a *Delivery) AddDocument(c *gin.Context) {
	req := new(AddDocumentRequest)
	if err := c.ShouldBind(req); err != nil {
		rest.ReError(c, http.StatusBadRequest, 0, err)
		return
	}

	if err := a.svc.AddDocuments(c, req.Texts); err != nil {
		rest.ReError(c, http.StatusInternalServerError, 0, err)
		return
	}

	rest.ReData(c, "")
}

func (a *Delivery) QueryDocument(c *gin.Context) {
	req := new(QueryDocumentsRequest)
	if err := c.ShouldBind(req); err != nil {
		rest.ReError(c, http.StatusBadRequest, 0, err)
		return
	}

	res, err := a.svc.QueryDocuments(c, req.Content)
	if err != nil {
		rest.ReError(c, http.StatusInternalServerError, 0, err)
		return
	}

	resp := new(QueryDocumentsResponse)
	resp.Content = res
	rest.ReData(c, resp)
}
