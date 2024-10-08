package aidoc

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/internal/infra/repository/vector/rag"
	"github.com/rubberduckkk/ducker/internal/service/aidoc"
	"github.com/rubberduckkk/ducker/pkg/ginhelper"
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
		cfg := config.Get().LLM.RedisVector
		ragRepo, err := rag.NewRedisRepo(
			rag.WithEmbedder(llms.Instance().Embedder),
			rag.WithScheme(cfg.Scheme),
			rag.WithHost(cfg.Host),
			rag.WithIndexName(cfg.Index),
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
		ginhelper.ReError(c, http.StatusBadRequest, 0, err)
		return
	}

	if err := a.svc.AddDocuments(c, req.Texts); err != nil {
		ginhelper.ReError(c, http.StatusInternalServerError, 0, err)
		return
	}

	ginhelper.ReData(c, "")
}

func (a *Delivery) QueryDocument(c *gin.Context) {
	req := new(QueryDocumentsRequest)
	if err := c.ShouldBind(req); err != nil {
		ginhelper.ReError(c, http.StatusBadRequest, 0, err)
		return
	}

	res, err := a.svc.QueryDocuments(c, req.Content)
	if err != nil {
		ginhelper.ReError(c, http.StatusInternalServerError, 0, err)
		return
	}

	resp := new(QueryDocumentsResponse)
	resp.Content = res
	ginhelper.ReData(c, resp)
}
