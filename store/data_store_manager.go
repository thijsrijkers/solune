package store

type DataStoreManager struct {
	stores map[string]*KeyValueStore
}

func NewDataStoreManager() *DataStoreManager {
	return &DataStoreManager{
		stores: make(map[string]*KeyValueStore),
	}
}

func (manager *DataStoreManager) AddStore(name string) {
	manager.stores[name] = NewKeyValueStore()
}

func (manager *DataStoreManager) GetStore(name string) (*KeyValueStore, bool) {
	store, exists := manager.stores[name]
	return store, exists
}
