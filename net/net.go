package net

import (
	"fmt"
	"net"
	"strconv"
)

func genErrInvalidIp(ip string) error {
	return fmt.Errorf("IP FORMAT INVALID: %v", ip)
}

func genErrInvalidPort(port string) error {
	return fmt.Errorf("PORT FORMAT INVALID: %v", port)
}


func GetAddrByIpPort(ip string, port int) (*net.TCPAddr, error) {
	if i := net.ParseIP(ip); i == nil || i.String() != ip {
		return nil, genErrInvalidIp(ip)
	}
	if port > 65535 || port < 1 {
		return nil, genErrInvalidPort(strconv.Itoa(port))
	}
	addrString := fmt.Sprintf("%v:%v", ip, port)
	return net.ResolveTCPAddr("tcp", addrString)
}

func GetIpPortByAddr(addr string) (string, int, error) {
	rAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return "", 0, err
	}
	return rAddr.IP.String(), rAddr.Port, nil
}

func ListenTcp(addr string) (*net.TCPListener, error) {
	rAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	ln, err := net.ListenTCP("tcp", rAddr)
	if err != nil {
		return nil, err
	}
	return ln, err
}
