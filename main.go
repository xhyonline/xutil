package main

import (
	"time"

	"github.com/xhyonline/xutil/xlog"
)

var logger = xlog.Get().Product("./test/a.log", true)

func main() {

	for {
		logger.Infof("你好世界")
		time.Sleep(time.Second)
	}
}
