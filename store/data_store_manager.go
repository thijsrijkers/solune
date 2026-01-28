package store

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"solune/filestore"
	"strconv"
	"strings"
)

type DataStoreManager struct {
	stores map[string]*KeyValueStore
}

func NewDataStoreManager() *DataStoreManager {
	manager := &DataStoreManager{
		stores: make(map[string]*KeyValueStore),
	}

	dbPath := filepath.Join("db")
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
		store, exists := manager.stores[storeName]
		if !exists || store == nil {
			log.Printf("Failed to initialize store: %s", storeName)
			continue
		}

		filePath := filepath.Join(dbPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v\n", filePath, err)
			continue
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || !strings.Contains(line, ",") {
				continue
			}

			parts := strings.SplitN(line, ",", 2)
			keyStr := strings.TrimSpace(parts[0])
			encodedValue := strings.TrimSpace(parts[1])

			valueBytes, err := base64.StdEncoding.DecodeString(encodedValue)
			if err != nil {
				log.Printf("Invalid base64 in %s for key %s: %v\n", filePath, keyStr, err)
				continue
			}

			keyInt, err := strconv.Atoi(keyStr)
			if err != nil {
				log.Printf("Invalid integer key %s in file %s\n", keyStr, filePath)
				continue
			}

			if err := store.Set(keyInt, string(valueBytes)); err != nil {
				log.Printf("Failed to load key %d into store %s: %v\n", keyInt, storeName, err)
			}
		}
	}

	return manager
}

func (manager *DataStoreManager) AddStore(name string) {
	fs, err := filestore.New(name)
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
	store, exists := manager.stores[name]
	if !exists {
		return false
	}

	if store.fileStore != nil {
		if err := store.fileStore.Close(); err != nil {
			log.Printf("Failed to close FileStore for %s: %v", name, err)
		}
	}

	delete(manager.stores, name)

	dbPath := filepath.Join("db")
	fileName := name + ".solstr"
	fullPath := filepath.Join(dbPath, fileName)

	if err := os.Remove(fullPath); err != nil {
		log.Printf("Failed to remove file %s: %v", fullPath, err)
	}

	return true
}
