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
	
	writer := bufio.NewWriter(conn)
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
			writeError(writer, err)
			return
		}

		result, err := s.execute(b.Instruction, b.Store, b.Key, b.Data)
		if err != nil {
			writeError(writer, err)
			return
		}
		
		writeResult(writer, result)
	}
}

func (s *Server) execute(action, storeName, key, data string) ([]map[string]interface{}, error) {
	switch action {
	case "get":
		return s.handleGet(storeName, key)
	case "set":
		return s.handleSet(storeName, key, data)
	case "delete":
		return s.handleDelete(storeName, key)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}

func (s *Server) handleGet(storeName string, key string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	if key != "" {
		uuidKey, err := uuid.Parse(key)
		if err != nil {
			return nil, fmt.Errorf("error parsing UUID: %s", err)
		}

		data, err := store.Get(uuidKey)
		if err != nil {
			return nil, err
		}
		return []map[string]interface{}{data}, nil
	}
	return store.GetAllData(), nil
}

func (s *Server) handleDelete(storeName string, key string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	if key != "" {
		uuidKey, err := uuid.Parse(key)
		if err != nil {
			return nil, fmt.Errorf("error parsing UUID: %s", err)
		}

		err = store.Delete(uuidKey)
		if err != nil {
			return nil, fmt.Errorf("failed to delete data: %s", err)
		}

		return []map[string]interface{}{{"status": 200}}, nil
	}

	removed := s.manager.RemoveStore(storeName)
	if removed {
		return []map[string]interface{}{{"status": 200}}, nil
	}

	return nil, fmt.Errorf("failed to remove store '%s'", storeName)
}



func (s *Server) handleSet(storeName string, key string, data string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		s.manager.AddStore(storeName)
		store, _ = s.manager.GetStore(storeName)
		return []map[string]interface{}{{"status": 200}}, nil
	}

	var parsedData map[string]interface{}
	if data != "" {
		err := json.Unmarshal([]byte(data), &parsedData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse data: %s", err)
		}
	}

	if key != "" {
		err := store.Update(key, parsedData)
		if err != nil {
			return nil, fmt.Errorf("failed to update data: %s", err)
		}
	} else {
		newUUID := uuid.New().String()
		err := store.Set(newUUID, parsedData)
		if err != nil {
			return nil, fmt.Errorf("failed to set data: %s", err)
		}
	}

	return []map[string]interface{}{{"status": 200}}, nil
}


func handleReadError(err error) {
	if err == io.EOF {
		return
	}
	fmt.Println("Read error:", err)
}

func writeError(writer *bufio.Writer, err error) {
	resp := map[string]interface{}{
		"error": err.Error(),
	}
	jsonData, _ := json.Marshal(resp)
	writer.WriteString(string(jsonData) + "\n")
	writer.Flush()
}

func writeResult(writer *bufio.Writer, result []map[string]interface{}) {
    if len(result) > 0 {
        for _, item := range result {
            jsonData, err := json.Marshal(item)
            if err != nil {
                writer.WriteString(`{"error": "failed to serialize data to JSON"}` + "\n")
                writer.Flush()
                fmt.Println("Error marshaling JSON:", err)
                return
            }
            dataToSend := string(jsonData) + "\n"
            writer.WriteString(dataToSend)
        }
    } else {
        dataToSend := `{"status": 404}` + "\n"
        writer.WriteString(dataToSend)
    }
    writer.Flush()
}
