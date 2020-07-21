package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Get 取得默认 logger 没有则创建一个
func Get(isDebug bool, path ...string) *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}
	switch {
	// 当且不是 Debug 状态下才会将日志追加到 log 文件中 ,讲道理 debug 的状态下直接输出在屏幕,谁还看日志啊......
	case len(path) != 0 && !isDebug:
		logfile, _ := os.OpenFile(path[0], os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		logger.SetOutput(logfile) //默认为os.stderr
	// 如果是 debug 或者是生产环境下没有设置日志的路径,也打在屏幕上
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
		})
		// 开启行号
		logger.SetReportCaller(true)
	}
	return logger
}
