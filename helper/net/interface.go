package net

import (
	"errors"
	"net"

	"github.com/hoquangnam45/pharmacy-common-go/util"
)

func FindBindInterfaceAddress(bindInterface string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, interf := range interfaces {
		if bindInterface != interf.Name {
			continue
		}

		addrs, err := interf.Addrs()
		if err != nil {
			return "", err
		}

		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.To4().String(), nil
				}
			}
		}
	}

	return "", errors.New("not found bind interface address")
}

func FindFirstNonLoopBackAddr() (*util.Pair[string, string], error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, interf := range interfaces {
		addrs, err := interf.Addrs()
		if err != nil {
			return nil, err
		}

		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return util.NewPair(ipnet.IP.To4().String(), interf.Name), nil
				}
			}
		}
	}

	return nil, errors.New("not found non loopback interface address")
}
