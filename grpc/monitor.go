package grpc

import (
	"net/http"

	"github.com/xhyonline/xutil/logger"
)

// pprofMonitor pprof 监控
func (s *grpcInstance) pprofMonitor() {
	go func() {
		if err := http.ListenAndServe(s.ip+":0", nil); err != nil {
			logger.Fatalf("pprof 服务启动失败")
		}
	}()
}
