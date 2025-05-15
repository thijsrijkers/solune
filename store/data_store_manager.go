package store

import (
	"log"
	"os"
	"path/filepath"
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
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".solstr") {
			continue
		}

		storeName := strings.TrimSuffix(name, ".solstr")

		manager.AddStore(storeName)
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
