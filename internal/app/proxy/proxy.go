package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloudflare/tableflip"
	"github.com/sirupsen/logrus"

	"github.com/rubberduckkk/ducker/internal/infra/config"
	"github.com/rubberduckkk/ducker/internal/infra/repository/file/account"
	"github.com/rubberduckkk/ducker/internal/service/proxy"
	"github.com/rubberduckkk/ducker/pkg/safe"
)

func Run(conf string) {
	config.Load(conf)

	if config.IsProd() {
		runGrace()
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Get().Port))
	if err != nil {
		panic(fmt.Sprintf("create listener failed: %v", err))
	}
	defer listener.Close()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	safe.Go(func() {
		defer wg.Done()
		runHTTPServer(ctx, listener)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	logrus.Infof("received shutdown signal. starting graceful shutdown")
	cancel()
	wg.Wait()
	logrus.Infof("graceful shutdown finished")
}

func runGrace() {
	upg, err := tableflip.New(tableflip.Options{PIDFile: "/var/run/ducker-proxy.pid"})
	if err != nil {
		panic(fmt.Sprintf("create tableflip Upgrader failed: %v", err))
	}
	defer upg.Stop()

	listener, err := upg.Listen("tcp", fmt.Sprintf(":%v", config.Get().Port))
	if err != nil {
		panic(fmt.Sprintf("create listener failed: %v", err))
	}
	defer listener.Close()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	safe.Go(func() {
		defer wg.Done()
		runHTTPServer(ctx, listener)
	})

	safe.Go(func() {
		sig := make(chan os.Signal, 1)
		// SIGTERM is signaled by k8s when it wants a pod to stop.
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		for range sig {
			logrus.Infof("received shutdown signal. starting graceful shutdown")
			cancel()
			if err = upg.Upgrade(); err != nil {
				logrus.WithError(err).Errorf("upgrade failed")
				return
			}
			logrus.Infof("tableflip upgrade success")
			return
		}
	})

	if err = upg.Ready(); err != nil {
		panic(fmt.Sprintf("call tableflip Ready failed: %v", err))
	}

	wg.Wait()
	<-upg.Exit()
	logrus.Infof("graceful shutdown finished")
}

func runHTTPServer(ctx context.Context, l net.Listener) {
	server := &http.Server{}

	svc := proxy.New(account.NewRepo(config.Get()))
	http.HandleFunc("/", svc.ProxyHTTP)

	safe.Go(func() {
		if err := server.Serve(l); err != nil {
			logrus.WithError(err).Errorf("serve http failed")
		}
	})

	<-ctx.Done()
	if err := server.Close(); err != nil {
		logrus.WithError(err).Errorf("close http server failed")
	}
}
