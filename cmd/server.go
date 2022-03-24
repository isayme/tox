package cmd

import (
	"io"
	"log"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/conf"
	"github.com/isayme/tox/socks5"
	"github.com/isayme/tox/tunnel"
	"github.com/isayme/tox/util"
)

func startServer() {
	config := conf.Get()

	formatTunnel, err := util.FormatURL(config.Tunnel)
	if err != nil {
		logger.Errorf("tunnel '%s' not valid format", config.Tunnel)
		return
	}
	config.Tunnel = formatTunnel

	ts, err := tunnel.NewServer(config.Tunnel, config.Password)
	if err != nil {
		logger.Errorw("new tunnel server fail", "err", err)
		return
	}

	logger.Infow("start listen", "addr", config.Tunnel)
	err = ts.ListenAndServeTLS(config.CertFile, config.KeyFile, func(rw io.ReadWriter) {
		request := socks5.NewRequest(rw)
		if err := request.Handle(); err != nil {
			logger.Errorw("socks5 fail", "err", err)
		}
	})

	log.Fatal(err)
}
