package rest

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/internal/infra/repository/vector/rag"
	"github.com/rubberduckkk/ducker/internal/service/aidoc"
	"github.com/rubberduckkk/ducker/pkg/llms"
)

type AIDocDelivery struct {
	svc aidoc.Service
}

var (
	r     *AIDocDelivery
	rOnce sync.Once
)

func AIDoc() *AIDocDelivery {
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
		r = &AIDocDelivery{svc: svc}
	})
	return r
}

func (a *AIDocDelivery) AddDocument(c *gin.Context) {
	req := new(AddDocumentRequest)
	if err := c.ShouldBind(req); err != nil {
		reError(c, http.StatusBadRequest, 0, err)
		return
	}

	if err := a.svc.AddDocuments(c, req.Texts); err != nil {
		reError(c, http.StatusInternalServerError, 0, err)
		return
	}

	reData(c, "")
}

func (a *AIDocDelivery) QueryDocument(c *gin.Context) {
	req := new(QueryDocumentsRequest)
	if err := c.ShouldBind(req); err != nil {
		reError(c, http.StatusBadRequest, 0, err)
		return
	}

	res, err := a.svc.QueryDocuments(c, req.Content)
	if err != nil {
		reError(c, http.StatusInternalServerError, 0, err)
		return
	}

	resp := new(QueryDocumentsResponse)
	resp.Content = res
	reData(c, resp)
}
