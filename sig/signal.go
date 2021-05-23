// sig 是接收信号通知的包

package sig

import (
	"os"
	"os/signal"
	"syscall"
)

type Func interface {
	GracefulClose()
}

// RegisterOnClose 注册关闭方法
func RegisterOnClose(f Func) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-c
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			f.GracefulClose()
		default:
			break
		}
	}()
}
