package conf

import (
	"sync"

	"github.com/isayme/go-config"
	"github.com/isayme/go-logger"
)

type Config struct {
	LogLevel      string `json:"log_level" yaml:"log_level"`
	Method        string `json:"method" yaml:"method"`
	Password      string `json:"password" yaml:"password"`
	Timeout       int    `json:"timeout" yaml:"timeout"`
	RemoteAddress string `json:"remote_address" yaml:"remote_address"`
	CertFile      string `json:"cert_file" yaml:"cert_file"`
	KeyFile       string `json:"key_file" yaml:"key_file"`
	LocalAddress  string `json:"local_address" yaml:"local_address"`
}

var once sync.Once
var globalConfig Config

func Get() *Config {
	config.Parse(&globalConfig)
	once.Do(func() {
		logger.SetLevel(globalConfig.LogLevel)
	})
	return &globalConfig
}
