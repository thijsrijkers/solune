package store

import (
	"log"
	"solune/filestore"
)

type DataStoreManager struct {
	stores map[string]*KeyValueStore
	port   string
}

func NewDataStoreManager(port string) *DataStoreManager {
	return &DataStoreManager{
		stores: make(map[string]*KeyValueStore),
		port:   port,
	}
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
