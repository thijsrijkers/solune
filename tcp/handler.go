package tcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			handleReadError(err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmd, err := ParseCommand(input)
		if err != nil {
			writeError(writer, err)
			continue
		}

		result, err := s.execute(cmd.Instruction, cmd.Store, cmd.Key, cmd.Data)
		if err != nil {
			writeError(writer, err)
			continue
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

func (s *Server) handleGet(storeName, key string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	if key == "" {
		allData := store.GetAllData()
		result := make([]map[string]interface{}, 0, len(allData))
		for k, v := range allData {
			result = append(result, map[string]interface{}{
				"key":   k,
				"value": v,
			})
		}
		return result, nil
	}

	keyInt, err := strconv.Atoi(key)
	if err != nil {
		return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
	}

	value, err := store.Get(keyInt)
	if err != nil {
		return nil, err
	}

	return []map[string]interface{}{
		{"key": keyInt, "value": value},
	}, nil
}

func (s *Server) handleDelete(storeName string, key string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		return nil, fmt.Errorf("store '%s' not found", storeName)
	}

	if key != "" {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
		}

		err = store.Delete(keyInt)
		if err != nil {
			return nil, fmt.Errorf("failed to delete key %d: %v", keyInt, err)
		}

		return []map[string]interface{}{{"status": 200}}, nil
	}

	if s.manager.RemoveStore(storeName) {
		return []map[string]interface{}{{"status": 200}}, nil
	}
	return nil, fmt.Errorf("failed to remove store '%s'", storeName)
}

func (s *Server) handleSet(storeName string, key string, data string) ([]map[string]interface{}, error) {
	store, exists := s.manager.GetStore(storeName)
	if !exists {
		s.manager.AddStore(storeName)
		store, _ = s.manager.GetStore(storeName)
	}

	if key != "" {
		keyInt, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid integer key '%s': %v", key, err)
		}

		err = store.Update(keyInt, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update key %d: %v", keyInt, err)
		}
	} else {
		all := store.GetAllData()
		newKey := len(all) + 1
		err := store.Set(newKey, data)
		if err != nil {
			return nil, fmt.Errorf("failed to set data: %v", err)
		}
	}

	return []map[string]interface{}{{"status": 200}}, nil
}

func handleReadError(err error) {
	if err != io.EOF {
		log.Println("Read error:", err)
	}
}

func writeError(writer *bufio.Writer, err error) {
	resp := map[string]interface{}{"error": err.Error()}
	jsonData, _ := json.Marshal(resp)
	writer.WriteString(string(jsonData) + "\n")
	writer.Flush()
}

func writeResult(writer *bufio.Writer, result []map[string]interface{}) {
	if len(result) == 0 {
		writer.WriteString(`{"status":404}` + "\n")
		writer.Flush()
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, item := range result {
		if err := enc.Encode(item); err != nil {
			log.Println("Error encoding item:", err)
			buf.WriteString(`{"error":"failed to serialize"}` + "\n")
		}
	}

	writer.Write(buf.Bytes())
	writer.Flush()
}
