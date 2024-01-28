package pgspy

import (
	"log"
	"net"
	"sync/atomic"
)

// ProxyConn - Proxy connection, piping data between proxy and remote.
type ProxyConn struct {
	laddr, raddr *net.TCPAddr
	lconn, rconn *net.TCPConn
	erred        bool
	errsig       chan bool
	connID       uint64
	msgID        uint64
	parser       *Parser
	collector    *QueriesCollector
	dataFilePath string
	lastQuery    *SQLQuery
}

// Pipe will move the bytes from the proxy to postgres and back
func (pc *ProxyConn) Pipe() {
	pc.collector = newQueriesCollector(pc.raddr.IP.String(), pc.dataFilePath)
	defer pc.collector.saveDataToFile()
	defer pc.lconn.Close()

	// connect to remote server
	rconn, err := net.DialTCP("tcp", nil, pc.raddr)
	if err != nil {
		log.Printf("Error: Failed to connect to Postgres: %s", err)
		return
	}
	pc.rconn = rconn
	defer pc.rconn.Close()

	// run parser
	go pc.parser.Parse()

	// proxying data
	go pc.handleRequestConnection(pc.lconn, pc.rconn)
	go pc.handleResponseConnection(pc.rconn, pc.lconn)

	// wait for close...
	<-pc.errsig
}

// Proxy.handleIncomingConnection
func (pc *ProxyConn) handleRequestConnection(src, dst *net.TCPConn) {
	// directional copy (64k buffer)
	buff := make([]byte, 0xffff)

	for {
		n, err := src.Read(buff)
		if err != nil {
			return
		}

		msgID := atomic.AddUint64(&pc.msgID, 1)
		parserBuff := make([]byte, n)
		copy(parserBuff, buff)
		go pc.sendMessageToParser(parserBuff, msgID, false)

		n, err = dst.Write(buff[:n])
		if err != nil {
			return
		}
	}
}

// Proxy.handleResponseConnection
func (pc *ProxyConn) handleResponseConnection(src, dst *net.TCPConn) {
	// directional copy (64k buffer)
	buff := make([]byte, 0xffff)

	for {
		n, err := src.Read(buff)
		if err != nil {
			return
		}

		msgID := atomic.AddUint64(&pc.msgID, 1)

		parserBuff := make([]byte, n)
		copy(parserBuff, buff[:n])
		go pc.sendMessageToParser(parserBuff, msgID, true)

		n, err = dst.Write(buff[:n])
		if err != nil {
			log.Printf("Error: Write failed '%s'\n", err)
			return
		}
	}
}

func (pc *ProxyConn) sendMessageToParser(buffer []byte, msgID uint64, outgoing bool) {
	wireMsg := WireMessage{
		Buff:     buffer,
		MsgID:    msgID,
		Outgoing: outgoing,
	}
	pc.parser.Incoming <- wireMsg
}

func (pc *ProxyConn) onMessage(msg PostgresMessage) {
	switch msg.TypeIdentifier {
	case QueryIncoming:
		pc.lastQuery = &SQLQuery{Query: string(msg.Payload)}
	case ErrorResponseOutgoing:
		if pc.lastQuery != nil {
			pc.lastQuery.IsSuccess = false
			pc.collector.addQuery(pc.lastQuery)
			pc.lastQuery = nil
		}
	case DataRowOutgoing, EmptyQueryResponseOutgoing,
		FunctionCallResponseOutgoing, NoticeResponseOutgoing, NoDataOutgoing:
		if pc.lastQuery != nil {
			pc.lastQuery.IsSuccess = true
			pc.collector.addQuery(pc.lastQuery)
			pc.lastQuery = nil
		}
	case CloseIncoming, TerminateIncoming:
		pc.collector.saveDataToFile()
	}
}

func (pc *ProxyConn) passMessagesToCallback(outgoing <-chan PostgresMessage) {
	for {
		msg := <-outgoing
		pc.onMessage(msg)
	}
}
