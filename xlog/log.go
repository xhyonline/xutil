package xlog

import (
	"os"
	"reflect"
	"strconv"
	"strings"

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
//	*gorm.DB.LogMode(true)
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
		break
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
func Get(isDebug bool, path ...string) *MyLogger {
	if logger == nil {
		logger = new(MyLogger)
		logger.Logger = logrus.New()
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
