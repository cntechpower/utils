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

func GetLocalIps() (localIps []*net.IPNet, err error) {
	localIps = make([]*net.IPNet, 0)
	addrS, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrS {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIps = append(localIps, ipNet)
			}
		}
	}
	return
}

func GetFirstLocalIp() (ip string, err error) {
	addrS, err := GetLocalIps()
	if err != nil {
		return
	}
	if len(addrS) < 1 {
		err = fmt.Errorf("can't find local ip")
		return
	}
	ip = addrS[0].IP.String()
	return
}
