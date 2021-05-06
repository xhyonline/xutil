package etcd

import (
	"fmt"
	"time"

	"context"

	"go.etcd.io/etcd/clientv3"
)

func New(address ...string) (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   address, // like 127.0.0.1:22379
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Status(timeout, address[0])
	if err != nil {
		return nil, fmt.Errorf("检查etcd状态失败: %v", err)
	}
	return client, nil
}
