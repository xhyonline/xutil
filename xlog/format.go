package xlog

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

type formatter struct {
}

type message struct {
	Date      string       `json:"date"`
	Level     logrus.Level `json:"level"`
	FilePath  string       `json:"file_path"`
	Message   string       `json:"message"`
	Func      string       `json:"func"`
	TimeStamp int64        `json:"time_stamp"`
}

func (m *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer = new(bytes.Buffer)
	date := entry.Time.Format("2006-01-02 15:04:05")
	pc, file, line, ok := runtime.Caller(7)
	var (
		f        = "未知方法名"
		filePath = "未知文件路径"
	)
	if ok {
		f = runtime.FuncForPC(pc).Name()
		filePath = file + ":" + strconv.Itoa(line)
	}
	body, _ := json.Marshal(&message{
		TimeStamp: entry.Time.Unix(),
		Date:      date,
		FilePath:  filePath,
		Message:   entry.Message,
		Level:     entry.Level,
		Func:      f,
	})
	buffer.WriteString(string(body) + "\n")
	return buffer.Bytes(), nil
}
