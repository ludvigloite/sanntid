package localip

import (
	"net"
	"strings"
)

var localIP string

//Sets up localIP.
//net.TCPAddr is a struct containing an IP, Port and Zone (IPv6 scoped addressing zone)

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}
