package pgspy

import (
	"fmt"
	"hash/fnv"
	"log"
	"net"
	"time"
)

// NewProxy - initializes a proxy
func NewProxy(postgresAddr, proxyAddr string, dataFile string) *Proxy {
	return &Proxy{
		PostgresAddr: postgresAddr,
		ProxyAddr:    proxyAddr,
		DataFilePath: dataFile,
		OnMessage:    func(_ PostgresMessage) {},
	}
}

type Callback func(get []byte, msgID uint64)

type Proxy struct {
	PostgresAddr string
	ProxyAddr    string
	DataFilePath string
	OnMessage    func(PostgresMessage)
}

func (p *Proxy) Start() {
	fmt.Printf("Server successfuly started\n proxy %s -> postgres: %s\n", p.ProxyAddr, p.PostgresAddr)

	postgresAddr := ResolvedAddress(p.PostgresAddr)
	proxyAddr := ResolvedAddress(p.ProxyAddr)

	listener, err := net.ListenTCP("tcp", proxyAddr)
	if err != nil {
		log.Fatalf("ListenTCP of %s error:%v", proxyAddr, err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Warning: Failed to accept connection '%s'\n", err)
			continue
		}

		incoming := make(chan WireMessage)
		outgoing := make(chan PostgresMessage)
		parser := &Parser{
			Incoming: incoming,
			Outgoing: outgoing,
		}

		proxyConn := &ProxyConn{
			lconn:        conn,
			laddr:        proxyAddr,
			raddr:        postgresAddr,
			erred:        false,
			errsig:       make(chan bool),
			connID:       generateUniqueId(),
			parser:       parser,
			dataFilePath: p.DataFilePath,
		}
		log.Printf("New connection %03d", proxyConn.connID)
		go proxyConn.Pipe()
		go proxyConn.passMessagesToCallback(outgoing)
	}
}
func generateUniqueId() uint64 {
	s, _ := time.Now().MarshalText()
	h := fnv.New32a()
	h.Write([]byte(s))
	return uint64(h.Sum32())
}
