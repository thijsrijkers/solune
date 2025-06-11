package tcprelay

import (
	"bufio"
	"log"
	"net"
	"strings"
	"fmt"
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
		return fmt.Errorf("failed to start server on port %s: %w", r.Port, err)
	}
	defer listener.Close()

	log.Printf("RelayNode listening on port %s\n", r.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v\n", err)
			continue
		}
		go r.handleConnection(conn)
	}
}

func (r *RelayNode) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	log.Printf("[Port %s] Client connected: %s\n", r.Port, clientAddr)
	defer func() {
		log.Printf("[Port %s] Client disconnected: %s\n", r.Port, clientAddr)
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		responses := r.relayMessage(msg)

		for _, response := range responses {
			_, err := fmt.Fprintf(conn, "%s\n", response)
			if err != nil {
				log.Printf("Failed to send response to %s: %v\n", clientAddr, err)
			}
		}
	}
}

func (r *RelayNode) relayMessage(msg string) []string {
	var responses []string
	for _, port := range r.PeerPorts {
		p := port
		conn, err := net.Dial("tcp", "localhost:"+p)
		if err != nil {
			log.Printf("Failed to connect to peer %s: %v\n", p, err)
			continue
		}

		_, err = fmt.Fprintf(conn, "%s\n", msg)
		if err != nil {
			log.Printf("Failed to send to peer %s: %v\n", p, err)
			conn.Close()
			continue
		}

		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err == nil {
			responses = append(responses, fmt.Sprintf(strings.TrimSpace(reply)))
		} else {
			log.Printf("No response from peer %s: %v\n", p, err)
		}
		conn.Close()
	}
	return responses
}

