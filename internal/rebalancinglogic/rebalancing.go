package rebalancinglogic

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type ShardData struct {
	Port int
	Stores map[string][]string // storeName -> keys
}

// Placeholder for communicating with existing shards
func communicateWithExistingShard(port int, store string, keys []string) {
	fmt.Printf("[existing-shard:%d] Store: %s, Keys: %v\n", port, store, keys)
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

