package utl

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
)

func GetHostname() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		msg := "ERR: Get local hostname failed!"
		log.Print(msg)
		return "", errors.New(msg)
	}

	return name, err
}

func GetIPAddrs() ([]string, error) {
	name, err := GetHostname()
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		msg := "ERR: Get local IP address failed"
		log.Print(msg)
		return nil, errors.New(msg)
	}

	return addrs, nil
}

func GetNetIP() (ip string, err error) {
	ips, err := GetIPAddrs()
	if err != nil {
		return
	}

	for i := range ips {
		if loc := strings.Index(ips[i], "127.0.0."); loc < 0 {
			ip = ips[i]
			break
		}
	}

	return
}
