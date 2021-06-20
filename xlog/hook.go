package xlog

import "github.com/sirupsen/logrus"

type Hook struct {
	Do    func(entry *logrus.Entry) error
	Level []logrus.Level
}

func (h *Hook) Levels() []logrus.Level {
	return h.Level
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	return h.Do(entry)
}
