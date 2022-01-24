package cmd

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/conf"
	"github.com/isayme/tox/tunnel"
	"github.com/isayme/tox/util"
)

func startLocal() {
	config := conf.Get()

	formatTunnel, err := util.FormatURL(config.Tunnel)
	if err != nil {
		logger.Errorf("tunnel '%s' not valid format", config.Tunnel)
		return
	}
	config.Tunnel = formatTunnel

	addr := config.LocalAddress
	logger.Infof("listen on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorw("Listen fail", "err", err)
		return
	}
	defer l.Close()

	tc, err := tunnel.NewClient(config.Tunnel, config.Password)
	if err != nil {
		logger.Errorw("new tunnel client fail", "err", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorw("l.Accept fail", "err", err)
			continue
		}

		go handleConnection(conn, tc)
	}
}

func handleConnection(conn net.Conn, tc tunnel.Client) {
	config := conf.Get()

	logger.Infow("new connection", "remoteAddr", conn.RemoteAddr().String())
	defer conn.Close()

	remote, err := tc.Connect(context.Background())
	if err != nil {
		logger.Errorw("connect tunnel server fail", "err", err)
		return
	}
	defer remote.Close()

	logger.Debug("connect tunnel server ok")

	conn = util.NewTimeoutConn(conn, time.Duration(config.Timeout)*time.Second)
	go io.Copy(remote, conn)
	io.Copy(conn, remote)
}
