package micro

import (
	"context"

	"go.etcd.io/etcd/clientv3"
)

// MicroServiceRegister 微服务注册实例
type MicroServiceRegister struct {
	// 客户端
	client *clientv3.Client
	// 前缀
	prefix string
}

// NewMicroServiceRegister 实例化注册器
func NewMicroServiceRegister(client *clientv3.Client, prefix string) *MicroServiceRegister {
	return &MicroServiceRegister{
		client: client,
		prefix: prefix,
	}
}

// Register 注册服务
func (s *MicroServiceRegister) Register(name string, node *Node) {
	s.client.Put(context.Background(), "")
}
