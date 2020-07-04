package xlog

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

// Get 取得默认 logger 没有则创建一个
func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

// Debug 将logger设置为debug模式
func Debug() {
	if logger == nil {
		logger = logrus.New()
	}
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
}
