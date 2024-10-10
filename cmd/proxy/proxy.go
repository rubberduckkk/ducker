package main

import (
	"flag"

	"github.com/rubberduckkk/ducker/internal/app/proxy"
)

var (
	configFile = flag.String("conf", "", "Path to config file")
)

func main() {
	flag.Parse()
	proxy.Run(*configFile)
}
