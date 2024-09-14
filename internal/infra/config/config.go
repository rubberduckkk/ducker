package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env    string `yaml:"env"`
	Port   uint16 `yaml:"port"`
	MainDB MySQL  `yaml:"main_db"`
	LLM    LLM    `yaml:"llm"`
}

var (
	config     *Config
	configOnce sync.Once
)

func Load(filename string) {
	configOnce.Do(func() {
		config = &Config{}
	})
	buf, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}
	if err = yaml.Unmarshal(buf, config); err != nil {
		panic(fmt.Sprintf("failed to unmarshal config file: %v", err))
	}
	logrus.WithField("config", config).Info("loaded config")
}

func Get() Config {
	if config == nil {
		return Config{}
	}
	return *config
}
