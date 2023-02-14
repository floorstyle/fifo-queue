package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/floorstyle/fifo-queue/util"
)

type Configuration struct {
	GinReleaseMode        bool
	HTTPPort              int
	RedisHost             string
	RedisDB               int
	HealthCheckMaxRetries int64
}

func (c *Configuration) ReadFile(path string) {
	file, err := os.Open(path)
	util.Try(err)
	decoder := json.NewDecoder(file)
	util.Try(decoder.Decode(&c))
}

func NewConfiguration() *Configuration {
	var configuration Configuration
	path := "config.json"

	configuration.ReadFile(DevModeConfiguration(path))
	return &configuration
}

func DevModeConfiguration(path string) string {
	if os.Getenv("DEVELOPMENT_MODE") == "true" {
		path = fmt.Sprint("dev.", path)
	}
	return filepath.Join(os.Getenv("CONFIG_PATH"), path)
}
