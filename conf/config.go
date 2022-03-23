package conf

import (
	"sync"

	"github.com/isayme/go-config"
	"github.com/isayme/go-logger"
)

type Config struct {
	LogLevel              string `json:"log_level" yaml:"log_level"`
	Password              string `json:"password" yaml:"password"`
	Timeout               int    `json:"timeout" yaml:"timeout"`
	Tunnel                string `json:"tunnel" yaml:"tunnel"`
	CertFile              string `json:"cert_file" yaml:"cert_file"`
	KeyFile               string `json:"key_file" yaml:"key_file"`
	TLSInsecureSkipVerify bool   `json:"tls_insecure_skip_verify" yaml:"tls_insecure_skip_verify"`
	LocalAddress          string `json:"local_address" yaml:"local_address"`
}

var once sync.Once
var globalConfig Config

func Get() *Config {
	config.Parse(&globalConfig)
	once.Do(func() {
		logger.SetLevel(globalConfig.LogLevel)
		logger.SetFormat("console")
	})
	return &globalConfig
}
