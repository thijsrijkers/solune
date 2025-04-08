package tcp

import (
    "bufio"
    "fmt"
	"io"
    "net"
    "strings"
    "solune/store"
)

type Server struct {
    manager *store.DataStoreManager
}

func NewServer(manager *store.DataStoreManager) *Server {
    return &Server{manager: manager}
}

func (s *Server) HandleClient(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                return
            }
            fmt.Println("Read error:", err)
            return
        }

        input = strings.TrimSpace(input)

		conn.Write([]byte(input + "\n"))
        // TODO: Create Input logic
    }
}
