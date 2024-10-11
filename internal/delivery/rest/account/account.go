package account

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/rubberduckkk/ducker/internal/infra/config"
	repo "github.com/rubberduckkk/ducker/internal/infra/repository/file/account"
	"github.com/rubberduckkk/ducker/internal/service/account"
	"github.com/rubberduckkk/ducker/pkg/ginhelper"
)

type Delivery struct {
	svc account.Service
}

var (
	d     *Delivery
	dOnce sync.Once
)

func Deliver() *Delivery {
	dOnce.Do(func() {
		r := repo.NewRepo(config.Get())
		d = &Delivery{svc: account.New(r)}
	})
	return d
}

func (d *Delivery) Auth(c *gin.Context) {
	req := new(AuthRequest)
	if err := c.ShouldBind(req); err != nil {
		ginhelper.ReError(c, http.StatusBadRequest, 0, err)
		return
	}

	if err := d.svc.Auth(req.Username, req.Password); err != nil {
		logrus.WithError(err).WithField("req", req).Errorf("auth failed")
		ginhelper.ReError(c, http.StatusUnauthorized, 0, err)
		return
	}

	ginhelper.ReData(c, "OK")
}
