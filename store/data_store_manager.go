package store

import (
	"log"
	"os"
	"path/filepath"
	"encoding/base64"
	"strings"
	"solune/filestore"
)

type DataStoreManager struct {
	stores map[string]*KeyValueStore
	port   string
}

func NewDataStoreManager(port string) *DataStoreManager {
	manager := &DataStoreManager{
		stores: make(map[string]*KeyValueStore),
		port:   port,
	}

	dbPath := filepath.Join("db", port)
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return manager
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".solstr") {
			continue
		}

		storeName := strings.TrimSuffix(file.Name(), ".solstr")
		manager.AddStore(storeName)

		store := manager.stores[storeName]
		filePath := filepath.Join(dbPath, file.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v\n", filePath, err)
			continue
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" || !strings.Contains(line, ",") {
				continue
			}

			parts := strings.SplitN(line, ",", 2)
			key := strings.TrimSpace(parts[0])
			encodedValue := strings.TrimSpace(parts[1])

			valueBytes, err := base64.StdEncoding.DecodeString(encodedValue)
			if err != nil {
				log.Printf("Invalid base64 in %s for key %s: %v\n", filePath, key, err)
				continue
			}

			store.cache[key] = valueBytes
			store.data[key] = valueBytes
		}
	}

	return manager
}

func (manager *DataStoreManager) AddStore(name string) {
	fs, err := filestore.New(name, manager.port)
	if err != nil {
		log.Printf("Failed to create filestore for %s: %v", name, err)
		return
	}
	manager.stores[name] = NewKeyValueStore(fs)
}

func (manager *DataStoreManager) GetStore(name string) (*KeyValueStore, bool) {
	store, exists := manager.stores[name]
	return store, exists
}

func (manager *DataStoreManager) RemoveStore(name string) bool {
	if _, exists := manager.stores[name]; exists {
		delete(manager.stores, name)
		return true
	}
	return false
}
