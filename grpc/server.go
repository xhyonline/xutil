package grpc

import (
	"net"
	"strconv"

	"github.com/xhyonline/xutil/metrics"

	"go.etcd.io/etcd/clientv3"

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
	// 服务名
	name string
	// 服务向 etcd 续租的租期
	lease int64
	// prometheus 监控网关
	prometheusGateWay string
}

// startCheck 启动前的检查
func (s *grpcInstance) startCheck() {
	if s.ip == "" {
		s.ip = internalIP()
	}
	if s.name == "" {
		logger.Fatalf("服务启动失败,请设定服务名")
	}
	if s.etcd == nil {
		logger.Fatalf("服务启动失败,无法向 etcd 注册服务")
	}
	if s.lease == 0 {
		s.lease = 10
	}
}

// WithPort 自定义端口,默认将会随机生成端口进行注册
func WithPort(port int) Option {
	return func(s *grpcInstance) {
		s.port = port
	}
}

// WithIP 自定义 IP 地址,默认为 eth0 网卡对应的 ipv4 地址
// Windows 下没有 eth0 网卡的,默认将会使用 127.0.0.1 作为 IP
func WithIP(ip string) Option {
	return func(s *grpcInstance) {
		s.ip = ip
	}
}

// WithAppName 服务名
func WithAppName(name string) Option {
	return func(s *grpcInstance) {
		s.name = name
	}
}

// WithETCD 必填项, etcd 客户端实例
func WithETCD(client *clientv3.Client) Option {
	return func(s *grpcInstance) {
		s.etcd = client
	}
}

// WithLease 自定义租约,默认为 10,
// 这意味着你的 grpc 服务每 10s 秒将会向 etcd 发送心跳包进行续租
func WithLease(lease int64) Option {
	return func(s *grpcInstance) {
		s.lease = lease
	}
}

// WithPrometheus 注册 prometheus 监控
func WithPrometheus(gateway string) Option {
	return func(s *grpcInstance) {
		s.prometheusGateWay = gateway
	}
}

// GracefulClose 优雅停止
func (s *grpcInstance) GracefulClose() {
	logger.Info("服务" + s.name + "接收到关闭通知")
	s.GracefulStop()
	logger.Info("服务" + s.name + "已优雅停止")
}

// run 启动
func (s *grpcInstance) run() {
	go func() {
		if err := s.Serve(s.listener); err != nil {
			logger.Fatalf("服务 %s 启动失败 %s", s.name, err)
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
	l, err := net.Listen("tcp", g.ip+":"+strconv.Itoa(g.port))
	if err != nil {
		logger.Fatalf("%s 服务监听失败", err)
	}
	if g.port == 0 {
		g.port = l.Addr().(*net.TCPAddr).Port
	}
	address := g.ip + ":" + strconv.Itoa(g.port)
	g.listener = l
	// 注册服务
	f(g.Server)
	g.run()
	ctx := sig.Get().RegisterClose(g)
	// 服务监控
	g.pprofMonitor()
	if g.prometheusGateWay != "" {
		metrics.Init(g.prometheusGateWay, g.name)
	}
	// 服务注册
	if err := micro.NewMicroServiceRegister(g.etcd, schema, g.lease).
		Register(g.name, &micro.Node{
			Host: g.ip,
			Port: strconv.Itoa(g.port),
		}); err != nil {
		logger.Fatalf("服务注册失败 %s", err)
	}
	logger.Info("服务"+g.name, "已启动,启动地址:"+address)
	<-ctx.Done()
}

// internalIP 获取内网 IP
func internalIP() string {
	var address = "127.0.0.1"
	addr, err := helper.IntranetAddress()
	if err != nil {
		logger.Fatalf("获取内网地址失败,服务停止 %s", err)
	}
	v := addr["eth0"]
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
