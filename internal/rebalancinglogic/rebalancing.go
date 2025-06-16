package rebalancinglogic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func executeGetOnShard(port int, store string, key string) (json.RawMessage, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect to port %d: %w", port, err)
	}
	defer conn.Close()

	buf := make([]byte, 0, 256)
	buf = append(buf, "instruction=get|store="...)
	buf = append(buf, store...)
	buf = append(buf, "|key="...)
	buf = append(buf, key...)
	buf = append(buf, '\n')


	if _, err := conn.Write(buf); err != nil {
		return nil, fmt.Errorf("write to port %d: %w", port, err)
	}

	readBuf := make([]byte, 4096)
	n, err := conn.Read(readBuf)
	if err != nil {
		return nil, fmt.Errorf("read from port %d: %w", port, err)
	}

	result := bytes.TrimSpace(readBuf[:n])
	return json.RawMessage(result), nil
}

func executeProcessStore(port int, store string, data json.RawMessage) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	delete(obj, "key")

	cleanData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect to port %d: %w", port, err)
	}
	defer conn.Close()

	var buf bytes.Buffer
	buf.WriteString("instruction=set|store=")
	buf.WriteString(store)
	buf.WriteString("|data=")
	buf.Write(cleanData)
	buf.WriteByte('\n')

	if _, err := conn.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write to port %d: %w", port, err)
	}

	readBuf := make([]byte, 1024)
	n, err := conn.Read(readBuf)
	if err != nil {
		return fmt.Errorf("read from port %d: %w", port, err)
	}

	response := strings.TrimSpace(string(readBuf[:n]))
	fmt.Printf("[shard:%d] store response: %s\n", port, response)

	return nil
}

func deleteKeyOnShard(port int, store, key string) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect to port %d: %w", port, err)
	}
	defer conn.Close()

	command := fmt.Sprintf("instruction=delete|store=%s|key=%s\n", store, key)

	if _, err := conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("write to port %d: %w", port, err)
	}

	readBuf := make([]byte, 1024)
	n, err := conn.Read(readBuf)
	if err != nil {
		return fmt.Errorf("read from port %d: %w", port, err)
	}

	response := strings.TrimSpace(string(readBuf[:n]))
	fmt.Printf("[shard:%d] delete response: %s\n", port, response)

	return nil
}

func createStoreOnShard(port int, store string) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect to port %d: %w", port, err)
	}
	defer conn.Close()

	command := fmt.Sprintf("instruction=set|store=%s\n", store)

	if _, err := conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("write to port %d: %w", port, err)
	}

	readBuf := make([]byte, 1024)
	n, err := conn.Read(readBuf)
	if err != nil {
		return fmt.Errorf("read from port %d: %w", port, err)
	}

	response := strings.TrimSpace(string(readBuf[:n]))
	fmt.Printf("[shard:%d] create store response: %s\n", port, response)

	return nil
}

func Run() error {
	dbDir := "./db"
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return fmt.Errorf("reading db directory: %w", err)
	}

	shardPorts := shardPortsFromEntries(entries)
	if len(shardPorts) == 0 {
		return fmt.Errorf("no valid shards found")
	}

	for _, sourcePort := range shardPorts {
		shardPath := filepath.Join(dbDir, strconv.Itoa(sourcePort))
		err := filepath.WalkDir(shardPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".solstr") {
				return nil
			}

			storeName := d.Name()[:len(d.Name())-len(".solstr")]

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Bytes()
				commaIdx := bytes.IndexByte(line, ',')
				if commaIdx == -1 {
					continue
				}

				key := strings.TrimSpace(string(line[:commaIdx]))

				targetPort := chooseTargetShard(key, shardPorts)
				if targetPort == sourcePort {
					continue
				}

				if !storeExistsInShard(dbDir, targetPort, storeName) {
					if err := createStoreOnShard(targetPort, storeName); err != nil {
						fmt.Printf("Error creating store '%s' on shard %d: %v\n", storeName, targetPort, err)
						continue
					}
				}

				response, err := executeGetOnShard(sourcePort, storeName, key)
				if err != nil {
					fmt.Printf("Error fetching key '%s': %v\n", key, err)
					continue
				}

				err = executeProcessStore(targetPort, storeName, response)
				if err != nil {
					fmt.Printf("Error processing key '%s' on destination shard: %v\n", key, err)
					continue
				}

				err = deleteKeyOnShard(sourcePort, storeName, key)
				if err != nil {
					fmt.Printf("Error deleting key '%s' on source shard: %v\n", key, err)
				}
			}
			return scanner.Err()
		})
		if err != nil {
			fmt.Printf("error processing shard %d: %v\n", sourcePort, err)
		}
	}

	return nil
}

func storeExistsInShard(dbDir string, shardPort int, store string) bool {
	storePath := filepath.Join(dbDir, strconv.Itoa(shardPort), store+".solstr")
	_, err := os.Stat(storePath)
	return err == nil
}

func shardPortsFromEntries(entries []os.DirEntry) []int {
	ports := make([]int, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		port, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		ports = append(ports, port)
	}
	return ports
}


func chooseTargetShard(key string, shardPorts []int) int {
	// Use a uint32 hash for faster arithmetic and avoid overflow
	var hash uint32 = 2166136261 // FNV-1a offset basis
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return shardPorts[int(hash)%len(shardPorts)]
}
