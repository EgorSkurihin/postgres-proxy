package pgspy

import (
	"log"
	"net"
)

// ResolvedAddress of host.
func ResolvedAddress(host string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal("ResolveTCPAddr of host:", err)
	}
	return addr
}
