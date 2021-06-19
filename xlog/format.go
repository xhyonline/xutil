package xlog

import (
	"bytes"
	"encoding/json"
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
	body, _ := json.Marshal(&message{
		TimeStamp: entry.Time.Unix(),
		Date:      date,
		FilePath:  entry.Caller.File + ":" + strconv.Itoa(entry.Caller.Line),
		Message:   entry.Message,
		Level:     entry.Level,
		Func:      entry.Caller.Function,
	})

	buffer.WriteString(string(body) + "\n")
	return buffer.Bytes(), nil
}
