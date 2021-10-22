package grpc

import (
	"net"
	"strconv"

	"go.etcd.io/etcd/clientv3"

	"github.com/xhyonline/micro-server-framework/component"
	"github.com/xhyonline/micro-server-framework/configs"
	"github.com/xhyonline/xutil/micro"
	"github.com/xhyonline/xutil/sig"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/xhyonline/xutil/helper"
	"github.com/xhyonline/xutil/logger"
	"google.golang.org/grpc"
)

type Option func(s *grpcInstance)

// grpcInstance GRPC 实例
type grpcInstance struct {
	*grpc.Server
	// 句柄
	listener net.Listener
	// etcd
	etcd *clientv3.Client
	// 端口
	port int
	// 默认 eth0 网卡的 ipv4 地址 ,如果没有网卡则为 127.0.0.1
	ip string
}

// startCheck 启动前的检查
func (s *grpcInstance) startCheck() {
	if s.ip == "" {
		s.ip = internalIP()
	}
	if s.port == 0 {
		logger.Fatalf("服务启动失败,无端口")
	}
	if s.etcd == nil {
		logger.Fatalf("服务启动失败,无法向 etcd 注册服务")
	}
}

func WithPort(port int) Option {
	return func(s *grpcInstance) {
		s.port = port
	}
}

func WithIP(ip string) Option {
	return func(s *grpcInstance) {
		s.ip = ip
	}
}

func WithETCD(client *clientv3.Client) Option {
	return func(s *grpcInstance) {
		s.etcd = client
	}
}

// GracefulClose 优雅停止
func (s *grpcInstance) GracefulClose() {
	logger.Info("服务" + configs.Name + "接收到关闭通知")
	s.GracefulStop()
	logger.Info("服务" + configs.Name + "已优雅停止")
}

// run 启动
func (s *grpcInstance) run() {
	go func() {
		if err := s.Serve(s.listener); err != nil {
			logger.Fatalf("服务 %s 启动失败 %s", configs.Name, err)
		}
	}()
}

// StartGRPCServer 启动 GRPC 服务
func StartGRPCServer(f func(server *grpc.Server), option ...Option) {
	g := &grpcInstance{Server: grpc.NewServer(registerMiddleware()...)}
	for _, optionFunc := range option {
		optionFunc(g)
	}
	// 启动前检查
	g.startCheck()
	address := g.ip + ":" + strconv.Itoa(g.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("%s 监听失败 %s", address, err)
	}
	g.listener = l
	// 注册服务
	f(g.Server)
	g.run()
	ctx := sig.Get().RegisterClose(g)
	// 服务监控
	g.pprofMonitor()
	// TODO promethus 监控注册
	// 服务注册
	if err := micro.NewMicroServiceRegister(component.Instance.ETCD, schema, 10).
		Register(configs.Name, &micro.Node{
			Host: g.ip,
			Port: strconv.Itoa(g.port),
		}); err != nil {
		logger.Fatalf("服务注册失败 %s", err)
	}
	logger.Info("服务"+configs.Name, "已启动,启动地址:"+address)
	<-ctx.Done()
}

// internalIP 获取内网 IP
func internalIP() string {
	var address = "127.0.0.1"
	addr, err := helper.IntranetAddress()
	if err != nil {
		logger.Fatalf("获取内网地址失败,服务停止 %s", err)
	}
	v, _ := addr["eth0"]
	var ip net.IP
	for _, item := range v {
		if ip = item.To4(); ip != nil {
			break
		}
	}
	if ip != nil {
		address = ip.String()
	} else {
		logger.Errorf("未发现 IPv4 地址,将使用 %s 替代", address)
	}

	return address
}

// registerMiddleware 注册中间键
func registerMiddleware() []grpc.ServerOption {
	return []grpc.ServerOption{
		// 处理 panic
		grpcmiddleware.WithUnaryServerChain(
			grpcrecovery.UnaryServerInterceptor(recoveryInterceptor()),
		),
		grpc.WriteBufferSize(0), grpc.ReadBufferSize(0),
	}
}
