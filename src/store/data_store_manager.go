package store

type DataStoreManager struct {
	stores map[string]*KeyValueStore
}

func NewKeyValueStoreManager() *DataStoreManager {
	return &DataStoreManager{
		stores: make(map[string]*KeyValueStore),
	}
}

func (manager *DataStoreManager) AddStore(name string, store *KeyValueStore) {
	manager.stores[name] = store
}

func (manager *DataStoreManager) GetStore(name string) (*KeyValueStore, bool) {
	store, exists := manager.stores[name]
	return store, exists
}
