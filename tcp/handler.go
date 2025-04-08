package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"encoding/json"
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

		// Parse command (assuming you have some logic to parse it into an instruction and store name)
		b, err := ParseCommand(input)
		if err != nil {
			conn.Write([]byte("Error: " + err.Error() + "\n"))
			return
		}

		// Call the private execute method
		result, err := s.execute(b.Instruction, b.Store)
		if err != nil {
			conn.Write([]byte("Error: " + err.Error() + "\n"))
			return
		}

		// Format the response and send it back to the client
		if len(result) > 0 {
			for _, item := range result {
				// Marshal the map to JSON format
				jsonData, err := json.Marshal(item)
				if err != nil {
					conn.Write([]byte("Error: failed to serialize data to JSON\n"))
					return
				}
				conn.Write([]byte(fmt.Sprintf("%s\n", jsonData)))
			}
		} else {
			conn.Write([]byte("404\n"))
		}
	}
}


func (s *Server) execute(action, storeName string) ([]map[string]interface{}, error) {
	// Retrieve the store by name
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	// Handle the action
	switch action {
	case "get":
		return store.GetAllData(), nil
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}
