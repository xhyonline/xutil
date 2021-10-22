package grpc

import (
	"github.com/xhyonline/xutil/etcd"
	"github.com/xhyonline/xutil/logger"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

// NewGRPCClient 启动一个 GRPC 客户端
func NewGRPCClient(name string, client *clientv3.Client) (*grpc.ClientConn, error) {
	r := etcd.NewResolver(name, client)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+":///"+name,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("客户端连接失败 %s", err)
	}
	return conn, nil
}
