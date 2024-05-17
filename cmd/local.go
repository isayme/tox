package cmd

import (
	"context"
	"net"
	"sync"
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
	defer conn.Close()

	config := conf.Get()

	tcpConn, _ := conn.(*net.TCPConn)
	conn = util.NewTimeoutConn(conn, time.Duration(config.Timeout)*time.Second)

	logger.Infow("new connection", "client", conn.RemoteAddr().String())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	remote, err := tc.Connect(ctx)
	if err != nil {
		logger.Errorw("connect tunnel server fail", "err", err)
		return
	}
	var once sync.Once
	defer func() {
		once.Do(func() {
			remote.Close()
		})
	}()

	logger.Debug("connect tunnel server ok")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		var err error
		var n int64
		n, err = util.CopyBuffer(remote, conn)
		logger.Debugw("copy from client end", "n", n, "err", err)
		remote.CloseWrite()
	}()

	go func() {
		defer wg.Done()

		var err error
		var n int64
		n, err = util.CopyBuffer(conn, remote)
		logger.Debugw("copy from remote end", "n", n, "err", err)
		tcpConn.CloseWrite()
	}()

	wg.Wait()

	logger.Debug("handle end")
}
