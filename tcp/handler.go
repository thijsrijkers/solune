package tcp

import (
	"bufio"
	"fmt"
	"net"
	"solune/store"
	"strings"
)

type Server struct {
	manager *store.DataStoreManager
}

func NewServer(manager *store.DataStoreManager) *Server {
	return &Server{manager: manager}
}

func (s *Server) HandleClient(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			HandleReadError(err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmd, err := ParseCommand(input)
		if err != nil {
			WriteError(writer, err)
			continue
		}

		result, err := s.execute(cmd.Instruction, cmd.Store, cmd.Key, cmd.Data)
		if err != nil {
			WriteError(writer, err)
			continue
		}

		WriteResult(writer, result)
	}
}

func (s *Server) execute(action, storeName, key, data string) ([]map[string]interface{}, error) {
	switch action {
	case "get":
		return s.HandleGet(storeName, key)
	case "set":
		return s.HandleSet(storeName, key, data)
	case "delete":
		return s.HandleDelete(storeName, key)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}
