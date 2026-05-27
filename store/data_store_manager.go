package store

import (
	"log"
	"os"
	"path/filepath"
	"solune/filestore"
	"strings"
	"sync"
)

type DataStoreManager struct {
	stores map[string]*KeyValueStore
	mu sync.RWMutex
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
	}

	return manager
}

func (manager *DataStoreManager) AddStore(name string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	fs, err := filestore.New(name)
	if err != nil {
		log.Printf("Failed to create filestore for %s: %v", name, err)
		return
	}
	manager.stores[name] = NewKeyValueStore(fs)
}

func (manager *DataStoreManager) GetStore(name string) (*KeyValueStore, bool) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	store, exists := manager.stores[name]
	return store, exists
}

func (manager *DataStoreManager) GetStores() []string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	names := make([]string, 0, len(manager.stores))
  for name := range manager.stores {
		names = append(names, name)
  }
	return names
}

func (manager *DataStoreManager) RemoveStore(name string) bool {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	
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
