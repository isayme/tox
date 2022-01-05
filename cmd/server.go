package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/isayme/go-logger"
	"github.com/isayme/go-toh2/conf"
	"github.com/isayme/go-toh2/middleware"
	"github.com/isayme/go-toh2/socks5"
	"github.com/posener/h2conn"
)

func startServer() {
	config := conf.Get()

	if middleware.NotExist(config.Method) {
		logger.Errorf("method '%s' not support", config.Method)
		return
	}

	URL, err := url.Parse(config.RemoteAddress)
	if err != nil {
		logger.Errorf("parse remove_address fail: %s", err.Error())
		return
	}

	http.HandleFunc("/toh2", func(rw http.ResponseWriter, req *http.Request) {
		conn, err := h2conn.Accept(rw, req)
		if err != nil {
			logger.Infof("failed creating connection from %s: %s", req.RemoteAddr, err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		logger.Info("new connection from %s", req.RemoteAddr)

		mw := middleware.Get(config.Method)
		wrapConn := mw(conn, config.Password)

		request := socks5.NewRequest(wrapConn)
		if err := request.Handle(); err != nil {
			logger.Errorw("socks5 fail", "err", err)
			rw.Write([]byte("welcome"))
		}
	})

	addr := fmt.Sprintf(":%s", URL.Port())
	logger.Infof("listen on %s", addr)
	err = http.ListenAndServeTLS(addr, "./testdata/server.pem", "./testdata/server.key", nil)
	log.Fatal(err)
}
