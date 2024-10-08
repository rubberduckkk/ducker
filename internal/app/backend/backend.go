package backend

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/rubberduckkk/ducker/internal/delivery/rest"
	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/pkg/llms"
	"github.com/rubberduckkk/ducker/pkg/mysql"
	"github.com/rubberduckkk/ducker/pkg/safe"
)

func Run(conf string) {
	config.Load(conf)
	if err := mysql.Init(config.Get().MainDB); err != nil {
		logrus.WithError(err).Fatal("init mysql database failed")
	}

	if err := llms.Init(config.Get().LLM); err != nil {
		logrus.WithError(err).Fatal("init llms failed")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Get().Port))
	if err != nil {
		panic(fmt.Sprintf("create listener failed: %v", err))
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	safe.Go(func() {
		defer wg.Done()
		runAPIServer(ctx, listener)
	})

	sig := make(chan os.Signal, 1)
	// SIGTERM is signaled by k8s when it wants a pod to stop.
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig
	logrus.Infof("received shutdown signal. starting graceful shutdown")
	cancel()
	wg.Wait()
	logrus.Infof("graceful shutdown finished")
}

func runAPIServer(ctx context.Context, listener net.Listener) {
	// default is gin.DebugMode
	if config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	rest.SetupGin(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server := &http.Server{
		Handler: r,
	}

	safe.Go(func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Errorf("server failed")
		}
	})

	<-ctx.Done()
	logrus.Infof("shutting down api server")
	// use a clean context to wait indefinitely
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Errorf("server shutdown failed")
	}
	logrus.Infof("api server shutdown complete")
}
