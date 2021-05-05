package helper

import "net"

// IntranetAddress 获取内网地址
func IntranetAddress() (map[string]net.Addr, error) {
	ips := make(map[string]net.Addr)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		if err != nil {
			return nil, err
		}
		for _, v := range addresses {
			ips[byName.Name] = v
		}
	}
	return ips, nil
}
