package xlog

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/sirupsen/logrus"
)

var logger *MyLogger

// MyLogger 它实现了 GORM 的日志追加方法
type MyLogger struct {
	*logrus.Logger // 继承库中的所有方法
}

// Print 构造 GORM 日志打印方法,能将完整的 sql 语句追加到日志中以便 debug
// 	使用示例:
// 	var log = xlog.Get(false, "./log.log")
//	*gorm.DB.SetLogger(log)
//  *gorm.DB.LogMode(true)
func (logger *MyLogger) Print(values ...interface{}) {
	if values[0] != "sql" {
		return
	}
	sqlString := values[3].(string)

	var list []interface{}
	if reflect.TypeOf(values[4]).Kind() == reflect.Slice {
		s := reflect.ValueOf(values[4])
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface())
		}
	}
	for _, v := range list {
		sqlString = strings.Replace(sqlString, "?", stringFromAssertionFloat(v), 1)
	}
	logger.Info(sqlString)
}

// XHook 钩子方法,可以到处触发
func (logger *MyLogger) XHook(do func(entry *logrus.Entry) error, level ...logrus.Level) {
	logger.AddHook(&Hook{
		Do:    do,
		Level: level,
	})
}

// Debug 模式
func (logger *MyLogger) Debug() *MyLogger {
	logger.SetOutput(os.Stdout)
	return logger
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

// stringFromAssertionFloat 断言浮动的字符串
func stringFromAssertionFloat(number interface{}) string {
	var numberString string
	switch floatOriginal := number.(type) {
	case float64:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case float32:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int32:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int64:
		numberString = strconv.FormatInt(floatOriginal, 10)
	case []uint8:
		numberString = string(floatOriginal)
	case string:
		numberString = "'" + floatOriginal + "'"
	case bool:
		numberString = strconv.FormatBool(floatOriginal)
	}
	return numberString
}

// Get 获取一个日志实例
// 参数:
// 		isDebug 是否为调试模式,调试模式日志只会打印在终端,如果想要配合追加日志路径,请填写为 false
// 		path 日志路径
func Get() *MyLogger {
	if logger == nil {
		logger = new(MyLogger)
		logger.Logger = logrus.New()
		// info 级别以下的都不输出
		logger.SetLevel(logrus.InfoLevel)
		// 开启行号
		logger.SetReportCaller(true)
		// 日志格式化
		logger.SetFormatter(&formatter{})
	}
	return logger
}
