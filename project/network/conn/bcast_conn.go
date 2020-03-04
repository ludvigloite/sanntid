// +build !windows

package conn

import (
	"net"
	"os"
	"syscall"
)

func DialBroadcastUDP(port int) net.PacketConn {
	s, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP) //sets up socket s???
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1) //reusing address for s
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1) //allow broadcast
	syscall.Bind(s, &syscall.SockaddrInet4{Port: port})							//binds s to port

	f := os.NewFile(uintptr(s), "")
	conn, _ := net.FilePacketConn(f) //returns copy of the packet network connection corresponding to the open file f
	f.Close()

	return conn
}
