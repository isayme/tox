package conf

import (
	"sync"

	"github.com/isayme/go-config"
	"github.com/isayme/go-logger"
)

type Config struct {
	Logger LoggerConfig `json:"logger" yaml:"logger"`

	// options for both client & server
	Tunnel   string `json:"tunnel" yaml:"tunnel"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`

	// server options
	CertFile string `json:"certFile" yaml:"certFile"`
	KeyFile  string `json:"keyFile" yaml:"keyFile"`

	// client options
	ConnectTimeout     int    `json:"connectTimeout" yaml:"connectTimeout"`
	LocalAddress       string `json:"localAddress" yaml:"localAddress"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`
}

type LoggerConfig struct {
	Level  string           `json:"level" yaml:"level"`
	Format logger.LogFormat `json:"format" yaml:"format"`
}

var once sync.Once
var globalConfig Config

func Get() *Config {
	config.Parse(&globalConfig)
	once.Do(func() {
		logger.SetLevel(globalConfig.Logger.Level)
		logger.SetFormat(globalConfig.Logger.Format)

		logger.Debugf("log with level: %s, format %s", globalConfig.Logger.Level, globalConfig.Logger.Format)
	})
	return &globalConfig
}

func (conf *Config) Default() {
	if conf.ConnectTimeout <= 0 {
		conf.ConnectTimeout = 3
	}

	if conf.LocalAddress == "" {
		conf.LocalAddress = ":1080"
	}
}
