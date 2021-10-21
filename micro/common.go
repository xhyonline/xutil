package micro

import (
	"fmt"
)

// Node 服务节点
type Node struct {
	// 该节点的主机地址
	Host string `json:"host"`
	// 端口
	Port string `json:"port"`
}

// Validate 验证节点
func (n *Node) Validate() error {
	if n.Port == "" || n.Host == "" {
		return fmt.Errorf("节点信息不正确")
	}
	return nil
}
