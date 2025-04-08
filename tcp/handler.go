package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"github.com/google/uuid"
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
			handleReadError(err)
			return
		}

		input = strings.TrimSpace(input)
		b, err := ParseCommand(input)
		if err != nil {
			writeError(conn, err)
			return
		}

		result, err := s.execute(b.Instruction, b.Store, b.Key, b.Data)
		if err != nil {
			writeError(conn, err)
			return
		}

		writeResult(conn, result)
	}
}

func (s *Server) execute(action, storeName, key, data string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	switch action {
	case "get":
		return s.handleGet(store, key)
	case "set":
		return s.handleSet(store, data)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}

func (s *Server) handleGet(store *store.KeyValueStore, key string) ([]map[string]interface{}, error) {
	if key != "" {
		uuidKey, err := uuid.Parse(key)
		if err != nil {
			return nil, fmt.Errorf("Error parsing UUID: %s", err)
		}

		data, err := store.Get(uuidKey)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{data}, nil
	}
	return store.GetAllData(), nil
}

func (s *Server) handleSet(store *store.KeyValueStore, data string) ([]map[string]interface{}, error) {
	var parsedData map[string]interface{}
	if data != "" {
		err := json.Unmarshal([]byte(data), &parsedData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse data: %s", err)
		}
	}

	newUUID := uuid.New()
	err := store.Set(newUUID, parsedData)
	if err != nil {
		return nil, fmt.Errorf("failed to set data: %s", err)
	}

	return []map[string]interface{}{{"status": 200}}, nil
}

func handleReadError(err error) {
	if err == io.EOF {
		return
	}
	fmt.Println("Read error:", err)
}

func writeError(conn net.Conn, err error) {
	conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func writeResult(conn net.Conn, result []map[string]interface{}) {
	if len(result) > 0 {
		for _, item := range result {
			jsonData, err := json.Marshal(item)
			if err != nil {
				conn.Write([]byte("Error: failed to serialize data to JSON\n"))
				return
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", jsonData)))
		}
	} else {
		conn.Write([]byte("{'status': 404}\n"))
	}
}
