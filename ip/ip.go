package ip

import (
	"fmt"
	"log"
	"net"
)

func GetWireGuardServerIP(cidr string) *net.IPNet {
	serverIP, serverSubnet := ParseCIDR(cidr)
	serverSubnet.IP = nextIP(serverIP.Mask(serverSubnet.Mask))
	return serverSubnet
}

func ParseCIDR(cidr string) (net.IP, *net.IPNet) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panic(err)
	}
	return ip, ipnet
}

func ParseIP(ip string) net.IP {
	netip, _ := ParseCIDR(fmt.Sprintf("%s/32", ip))
	return netip
}

func nextIP(ip net.IP) net.IP {
	next := make([]byte, len(ip))
	copy(next, ip)
	for j := len(next) - 1; j >= 0; j-- {
		next[j]++
		if next[j] > 0 {
			break
		}
	}
	return next
}
