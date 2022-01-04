package cmd

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-toh2/conf"
	"github.com/isayme/go-toh2/middleware"
	"github.com/isayme/go-toh2/util"
)

func startLocal() {
	config := conf.Get()

	if middleware.NotExist(config.Method) {
		logger.Panicf("method '%s' not support", config.Method)
		return
	}

	addr := config.LocalAddress
	logger.Infof("listen on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorw("Listen fail", "err", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorw("l.Accept fail", "err", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	config := conf.Get()

	logger.Infow("new connection", "remoteAddr", conn.RemoteAddr().String())
	defer conn.Close()

	remote, resp, err := util.H2Client.Connect(context.Background(), config.RemoteAddress)
	if err != nil {
		logger.Errorw("client.Connect fail", "err", err)
		return
	}
	defer remote.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("bad status code: %d", resp.StatusCode)
		return
	}

	logger.Info("connect http2 server ok")

	md := middleware.Get(config.Method)
	wrapRemote := md(remote, config.Password)

	go io.Copy(wrapRemote, conn)
	io.Copy(conn, wrapRemote)
}
