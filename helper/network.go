package helper

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// IntranetAddress 获取内网地址
func IntranetAddress() (map[string][]net.IP, error) {
	ips := make(map[string][]net.IP)

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
			ipString := strings.Split(v.String(), "/")[0]
			if ip := net.ParseIP(ipString); ip != nil {
				ips[byName.Name] = append(ips[byName.Name], ip)
				continue
			}
			return nil, fmt.Errorf("ip 地址解析失败 %s", v.String())
		}
	}
	return ips, nil
}

// PublicNetAddress 获取公网地址
// 请保证机器能通外网
func PublicNetAddress() (net.IP, error) {
	body, err := HTTPRequest("http://ip.dhcp.cn/?ip", "GET", http.Header{}, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(string(body)); ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("公网 IP 解析失败 %s", string(body))
}
