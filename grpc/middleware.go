package grpc

import (
	"fmt"
	"runtime"

	"github.com/xhyonline/xutil/logger"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func recoveryInterceptor() grpcrecovery.Option {
	return grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
		buf := make([]byte, 4096)
		// 抛出服务端 调用栈的轨迹
		// runtime.stack 详情参考资料
		// https://colobu.com/2016/12/21/how-to-dump-goroutine-stack-traces/
		num := runtime.Stack(buf, false)
		msg := fmt.Sprintf("[grpc_panic_recovery]: %v %s", p, string(buf[:num]))
		logger.Errorf(msg)
		// p 就是捕获 panic 中的内容
		return status.Errorf(codes.Unknown, "%s", msg)
	})
}
