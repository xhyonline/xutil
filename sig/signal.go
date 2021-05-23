// sig 是接收信号通知的包

package sig

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once

var server *signalInstance

// Func 所有方法必须实现优雅退出,已保证退出后所做的事
type Func interface {
	GracefulClose()
}

// signalInstance 是一个信号抽象
type signalInstance struct {
	// 方法集合
	funcArr []Func
	// 退出信号
	closed chan struct{}
	// 系统信号
	signalChan chan os.Signal
}

// GracefulClose 注册关闭方法
func (s *signalInstance) RegisterClose(f Func) <-chan struct{} {
	server.funcArr = append(server.funcArr, f)
	return server.closed
}

// Get 获取单例
func Get() *signalInstance {
	once.Do(func() {
		server = &signalInstance{
			funcArr:    make([]Func, 0),
			closed:     make(chan struct{}, 1),
			signalChan: make(chan os.Signal),
		}
		go func() {
			signal.Notify(server.signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			sig := <-server.signalChan
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				for _, f := range server.funcArr {
					f.GracefulClose()
				}
				server.closed <- struct{}{}
			}
		}()
	})
	return server
}
