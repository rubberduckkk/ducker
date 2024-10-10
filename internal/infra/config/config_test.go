package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	Load("../../../config/config.dev.yaml")
}
