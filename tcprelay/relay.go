package tcprelay

import (
	"io"
	"log"
	"net"
	"time"
)

type RelayNode struct {
	Port      string
	PeerPorts []string
}

func NewRelayNode(port string, peerPorts []string) *RelayNode {
	return &RelayNode{
		Port:      port,
		PeerPorts: peerPorts,
	}
}

func (r *RelayNode) Start() error {
	listener, err := net.Listen("tcp", ":"+r.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("RelayNode listening on port %s\n", r.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go r.handleConnection(conn)
	}
}

func (r *RelayNode) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	clientAddr := clientConn.RemoteAddr().String()
	log.Printf("[OPEN] Client connected on port %s: %s\n", r.Port, clientAddr)

	if tcpConn, ok := clientConn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	if len(r.PeerPorts) == 0 {
		log.Printf("[ERROR] No peer ports available for port %s\n", r.Port)
		return
	}
	peerPort := r.PeerPorts[0]

	peerConn, err := net.Dial("tcp", "localhost:"+peerPort)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to peer %s: %v\n", peerPort, err)
		return
	}
	defer peerConn.Close()

	if tcpPeerConn, ok := peerConn.(*net.TCPConn); ok {
		tcpPeerConn.SetNoDelay(true)
		tcpPeerConn.SetKeepAlive(true)
		tcpPeerConn.SetKeepAlivePeriod(30 * time.Second)
	}

	bufSize := 64 * 1024
	done := make(chan struct{}, 2)

	// Client -> Peer
	go func() {
		_, err := io.CopyBuffer(peerConn, clientConn, make([]byte, bufSize))
		if err != nil {
			log.Printf("[ERROR] Copy client -> peer failed: %v\n", err)
		}
		done <- struct{}{}
	}()

	// Peer -> Client
	go func() {
		_, err := io.CopyBuffer(clientConn, peerConn, make([]byte, bufSize))
		if err != nil {
			log.Printf("[ERROR] Copy peer -> client failed: %v\n", err)
		}
		done <- struct{}{}
	}()

	<-done
	log.Printf("[CLOSE] Client disconnected on port %s: %s\n", r.Port, clientAddr)
}
