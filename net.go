package util

import (
	"errors"
	"net"
	"strings"
)

var (
	errorAddressFormat = errors.New("error address format")
)

// ResolveAddress 解析 tcp/udp:ip:port 的地址
func ResolveAddress(addr string) (net.Addr, error) {
	// network
	i := strings.IndexByte(addr, ':')
	if i < 0 {
		return nil, errorAddressFormat
	}
	return ResolveAddr(addr[:i], addr[i+1:])
}

// ResolveAddr 用标准库解析 network ip:port 的地址
func ResolveAddr(network, address string) (net.Addr, error) {
	network = strings.ToLower(network)
	// 解析
	switch network {
	case "tcp":
		return net.ResolveTCPAddr(network, address)
	default:
		return net.ResolveUDPAddr(network, address)
	}
}

// IPAddrOfFirstInterface 获取第一个网卡地址
func IPAddrOfFirstInterface() (string, error) {
	netifs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, netif := range netifs {
		// 关闭的
		if netif.Flags&net.FlagUp != 1 {
			continue
		}
		addrs, err := netif.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if _a, ok := addr.(*net.IPNet); ok {
				if _a.IP.To4() != nil || !_a.IP.IsLoopback() {
					return _a.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}
