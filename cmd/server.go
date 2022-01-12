package cmd

import (
	"io"
	"log"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/conf"
	"github.com/isayme/tox/middleware"
	"github.com/isayme/tox/socks5"
	"github.com/isayme/tox/tunnel"
)

func startServer() {
	config := conf.Get()

	if middleware.NotExist(config.Method) {
		logger.Errorf("method '%s' not support", config.Method)
		return
	}

	ts, err := tunnel.NewServer(config.Tunnel)
	if err != nil {
		logger.Errorw("new tunnel server fail", "err", err)
		return
	}

	err = ts.ListenAndServeTLS(config.CertFile, config.KeyFile, func(rw io.ReadWriter) {
		mw := middleware.Get(config.Method)
		wrapConn := mw(rw, config.Password)

		request := socks5.NewRequest(wrapConn)
		if err := request.Handle(); err != nil {
			logger.Errorw("socks5 fail", "err", err)
			rw.Write([]byte("welcome"))
		}
	})

	log.Fatal(err)
}
