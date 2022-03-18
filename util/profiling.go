package util

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/isayme/go-logger"
)

func EnableProfiling(enabled bool) {
	if enabled {
		go func() {
			addr := "0.0.0.0:6060"
			logger.Infof("start profiling %s", addr)
			http.ListenAndServe(addr, nil)
		}()
	}
}
