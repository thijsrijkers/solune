package rebalancinglogic

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"encoding/json"
	"strconv"
	"strings"
	"net"
	"sync"
)

type ShardData struct {
	Port int
	Stores map[string][]string // storeName -> keys
}

func communicateWithExistingShard(port int, store string, keys []string) []json.RawMessage {
	address := fmt.Sprintf("127.0.0.1:%d", port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Failed to connect to shard at %s: %v\n", address, err)
		return nil
	}
	defer conn.Close()

	fmt.Printf("[existing-shard:%d] Store: %s, Keys: %v\n", port, store, keys)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var results []json.RawMessage

	for _, key := range keys {
		command := key + "\n"

		_, err := writer.WriteString(command)
		if err != nil {
			fmt.Printf("Failed to write key '%s': %v\n", key, err)
			continue
		}
		writer.Flush()

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read response for key '%s': %v\n", key, err)
			continue
		}
		line = strings.TrimSpace(line)

		var maybeStatus map[string]interface{}
		if err := json.Unmarshal([]byte(line), &maybeStatus); err == nil {
			if len(maybeStatus) == 1 {
				if _, isStatus := maybeStatus["status"]; isStatus {
					fmt.Printf("Key '%s' returned status-only response, skipping.\n", key)
					continue
				}
			}
		}

		var raw json.RawMessage
		if err := json.Unmarshal([]byte(line), &raw); err == nil {
			results = append(results, raw)
		} else {
			fmt.Printf("Failed to parse JSON for key '%s': %v\n", key, err)
		}
	}

	return results
}

// Placeholder for sending data to new shard
func sendToNewShard(port int, store string, key string, value string) {
	fmt.Printf("[new-shard] Store: %s, Key: %s, Value: %s\n", store, key, value)
}

func Run(pids []int) error {
	dbDir := "./db"
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return fmt.Errorf("reading db directory: %w", err)
	}

	var wg sync.WaitGroup
	shardChan := make(chan ShardData)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		port, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue // skip non-numeric dirs
		}

		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			shardPath := filepath.Join(dbDir, fmt.Sprintf("%d", port))
			stores := make(map[string][]string)

			err := filepath.WalkDir(shardPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".solstr") {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				storeName := strings.TrimSuffix(d.Name(), ".solstr")
				var keys []string

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					parts := strings.SplitN(scanner.Text(), ",", 2)
					if len(parts) != 2 {
						continue
					}
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					keys = append(keys, key)

					// Communicate with new shard (placeholder)
					sendToNewShard(port, storeName, key, value)
				}
				if err := scanner.Err(); err != nil {
					return err
				}

				// Communicate with existing shard (placeholder)
				communicateWithExistingShard(port, storeName, keys)
				stores[storeName] = keys
				return nil
			})

			if err != nil {
				fmt.Printf("error walking shard %d: %v\n", port, err)
				return
			}

			shardChan <- ShardData{Port: port, Stores: stores}
		}(port)
	}

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(shardChan)
	}()

	// Process all shard results
	for shard := range shardChan {
		fmt.Printf("Finished shard %d with stores: %v\n", shard.Port, shard.Stores)
	}

	return nil
}

