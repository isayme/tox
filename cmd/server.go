package cmd

import (
	"time"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/conf"
	"github.com/isayme/tox/socks5"
	"github.com/isayme/tox/tunnel"
	"github.com/isayme/tox/util"
)

func startServer() {
	config := conf.Get()
	config.Default()

	formatTunnel, err := util.FormatURL(config.Tunnel)
	if err != nil {
		logger.Errorf("tunnel '%s' not valid format", config.Tunnel)
		return
	}
	config.Tunnel = formatTunnel

	if config.Password == "" {
		logger.Errorf("password is required")
		return
	}

	options := []util.ToxOption{
		util.WithTunnel(config.Tunnel),
		util.WithPassword(config.Password),
		util.WithTimeout(time.Second * time.Duration(config.Timeout)),
		util.WithCertFile(config.CertFile),
		util.WithKeyFile(config.KeyFile),
		util.WithConnectTimeout(time.Second * time.Duration(config.ConnectTimeout)),
	}

	ts, err := tunnel.NewServer(util.ToToxOptions(options))
	if err != nil {
		logger.Errorw("new tunnel server fail", "err", err)
		return
	}

	logger.Infow("start listen", "addr", config.Tunnel)
	err = ts.ListenAndServe(func(conn util.ToxConn) {
		/**
		 * return when server will not send data anymore
		 */
		request := socks5.NewRequest(config, conn)
		if err := request.Handle(); err != nil {
			logger.Errorw("socks5 fail", "err", err)
		}
	})
	if err != nil {
		logger.Errorf("listen fail %v", err)
	}
}
