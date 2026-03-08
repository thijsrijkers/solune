package tcp

import (
	"errors"
	"log"
	"net"
	"solune/store"
	"sync"
)

func StartServer(port string, manager *store.DataStoreManager) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Println("Listening on port", port)

	server := NewServer(manager)
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Listener was closed, exit cleanly
			if isClosedErr(err) {
				break
			}
			log.Println("Connection error:", err)
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			server.HandleClient(c)
		}(conn)
	}

	wg.Wait() // drain active connections before returning
}

func isClosedErr(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return netErr.Err.Error() == "use of closed network connection"
	}
	return false
}
