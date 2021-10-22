package grpc

import (
	"net/http"
	"os"

	"github.com/xhyonline/xutil/logger"
)

// pprofMonitor pprof 监控
func pprofMonitor() {
	go func() {
		if err := http.ListenAndServe(internalIP()+":0", nil); err != nil {
			logger.Errorf("pprof 服务启动失败")
			os.Exit(1)
		}
	}()
}

// TODO prometheus
