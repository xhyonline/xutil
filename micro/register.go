package micro

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/xhyonline/xutil/sig"

	"go.etcd.io/etcd/clientv3"
)

// MicroServiceRegister 微服务注册实例
type MicroServiceRegister struct {
	// 客户端
	client *clientv3.Client
	// 前缀
	prefix string
	// 租约时间
	lease int64
	// 全员信号
	ctx context.Context
	// 取消方法
	cancelFunc context.CancelFunc
}

// NewMicroServiceRegister 实例化注册器
func NewMicroServiceRegister(client *clientv3.Client, prefix string, lease int64) *MicroServiceRegister {
	if lease <= 0 {
		lease = 10
	}
	ctx, cancel := context.WithCancel(context.Background())
	register := &MicroServiceRegister{
		client:     client,
		prefix:     "/" + strings.Trim(prefix, "/") + "/",
		lease:      lease,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	// 优雅退出
	sig.Get().RegisterClose(register)
	return register
}

// GracefulClose 优雅退出
func (s *MicroServiceRegister) GracefulClose() {
	logger.Info("服务发现组件正在准备做退出的清理工作")
	s.cancelFunc()
	time.Sleep(time.Second * 5)
	defer s.client.Close()
	logger.Info("服务发现组件清理工作已完成")
}

// Register 注册服务
func (s *MicroServiceRegister) Register(name string, node *Node) error {
	v, err := json.Marshal(node)
	if err != nil {
		return err
	}

	resp, err := s.client.Grant(context.Background(), s.lease)
	if err != nil {
		return err
	}

	key := s.prefix + name + "/" + node.Host + ":" + node.Port
	if _, err = s.client.Put(context.Background(), key, string(v), clientv3.WithLease(resp.ID)); err != nil {
		return err
	}

	leaseChan, err := s.client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}

	go func(serverName string, leaseID clientv3.LeaseID, node *Node, ctx context.Context) {
		for {
			select {
			case leaseKeepResp := <-leaseChan:
				logger.Infof("服务名:%s 地址: %s 续租成功 %+v", serverName, node.Host+":"+node.Port, leaseKeepResp)
			case <-ctx.Done():
				if _, err := s.client.Revoke(context.Background(), leaseID); err != nil {
					logger.Errorf("撤销租约发生错误 %s", err)
				}
				logger.Infof("服务名:%s 地址: %s 已安全退出", serverName, node.Host+":"+node.Port)
				return
			}
		}
	}(name, resp.ID, node, s.ctx)

	return nil
}
