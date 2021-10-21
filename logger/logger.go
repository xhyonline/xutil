// logger 工具库作为 xlog 的封装
package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/xhyonline/xutil/xlog"
)

var (
	instance *xlog.MyLogger
)

func init() {
	instance = xlog.Get()
	// 默认 info 级别以下不打印
	instance.SetLevel(logrus.InfoLevel)
}

func SetLoggerProduct(path string) {
	instance = instance.Product(path, true)
}

func SetHook(do func(entry *logrus.Entry) error, level ...logrus.Level) {
	instance.XHook(do, level...)
}

func SetLevel(level logrus.Level) {
	instance.SetLevel(level)
}

func Infof(format string, args ...interface{}) {
	instance.Infof(format, args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Errorf(format string, args ...interface{}) {
	instance.Errorf(format, args...)
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Warnf(format string, args ...interface{}) {
	instance.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	instance.Debugf(format, args...)
}

func Panic(args ...interface{}) {
	instance.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	instance.Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	instance.Fatalf(format, args...)
}
