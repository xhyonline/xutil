package xlog

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/sirupsen/logrus"
)

// MyLogger 它实现了 GORM 的日志追加方法
type MyLogger struct {
	*logrus.Logger // 继承库中的所有方法
}

// XHook 钩子方法,可以到处触发
func (logger *MyLogger) XHook(do func(entry *logrus.Entry) error, level ...logrus.Level) {
	logger.AddHook(&Hook{
		Do:    do,
		Level: level,
	})
}

// Product 生产模式
// 参数:split 是否开启日志分割
func (logger *MyLogger) Product(path string, split bool) *MyLogger {
	if !split {
		logfile, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		logger.SetOutput(logfile) //默认为os.stderr
		return logger
	}
	writer, _ := rotatelogs.New(
		path+"_%Y%m%d%H",
		// WithLinkName, // 为最新的日志建立软连接
		// WithRotationCount //设置文件清理前最多保存的个数
		// WithMaxAge 和 WithRotationCount二者只能设置一个
		// 设置文件清理前的最长保存时间,一天后自动删除 (注:最小单位分钟)
		rotatelogs.WithMaxAge(time.Hour*24),
		// 设置日志分割的时间，每一个小时分割一次
		rotatelogs.WithRotationTime(time.Hour),
	)
	logger.SetOutput(writer)
	return logger
}

// Get 获取一个日志实例
func Get() *MyLogger {
	logger := new(MyLogger)
	logger.Logger = logrus.New()
	// 开启行号
	logger.SetReportCaller(true)
	// 日志格式化
	logger.SetFormatter(&formatter{})
	// 标准输出
	logger.SetOutput(os.Stdout)
	return logger
}
