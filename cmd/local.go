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

	addr := config.LocalAddress

	options := []util.ToxOption{
		util.WithTunnel(config.Tunnel),
		util.WithPassword(config.Password),
		util.WithTimeout(time.Second * time.Duration(config.Timeout)),
		util.WithLocalAddress(addr),
		util.WithConnectTimeout(time.Second * time.Duration(config.ConnectTimeout)),
		util.WithInsecureSkipVerify(config.InsecureSkipVerify),
	}

	logger.Infof("listen on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Errorw("Listen fail", "err", err)
		return
	}
	defer l.Close()

	tc, err := tunnel.NewClient(util.ToToxOptions(options))
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

		go handleConnection(config, conn, tc)
	}
}

func handleConnection(config *conf.Config, conn net.Conn, tc tunnel.Client) {
	defer conn.Close()

	tcpConn, _ := conn.(*net.TCPConn)
	conn = util.NewTimeoutConn(conn, time.Duration(config.Timeout)*time.Second)

	logger.Infow("new connection", "client", conn.RemoteAddr().String())

	remote, err := tc.Connect(context.Background())
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
