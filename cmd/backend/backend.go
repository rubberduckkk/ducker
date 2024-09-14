package main

import (
	"flag"

	"github.com/rubberduckkk/ducker/internal/app/backend"
)

var (
	configFile = flag.String("conf", "", "Path to config file")
)

func main() {
	flag.Parse()
	backend.Run(*configFile)
}
