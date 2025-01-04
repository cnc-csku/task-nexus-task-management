package network

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("No Local IP Found")
}
