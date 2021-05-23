package main

import (
	"fmt"
	"time"

	"github.com/xhyonline/xutil/sig"
)

type c struct {
}

func (s *c) GracefulClose() {

	fmt.Println("测试关闭")
	time.Sleep(time.Second * 10)
	fmt.Println("关闭结束")
}

func main() {
	s := new(c)
	sig.RegisterOnClose(s)
	select {}
}
