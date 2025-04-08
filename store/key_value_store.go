package store

import (
	"fmt"
	"sync"
)

type KeyValueStore struct {
	data   map[interface{}]map[string]interface{}
	cache  map[interface{}]map[string]interface{}
	mutex  sync.RWMutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data:  make(map[interface{}]map[string]interface{}),
		cache: make(map[interface{}]map[string]interface{}),
	}
}

func (store *KeyValueStore) Set(key interface{}, value map[string]interface{}) error {
	// Lock the store for writes
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.cache[key] = value
	
	store.data[key] = value
	return nil
}

func (store *KeyValueStore) Get(key interface{}) (map[string]interface{}, error) {
	if value, found := store.cache[key]; found {
		return value, nil
	}

	// Lock the store for reads
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if value, exists := store.data[key]; exists {
		store.cache[key] = value
		return value, nil
	}

	return nil, &KeyNotFoundError{Key: key}
}

func (store *KeyValueStore) GetAllData() []map[string]interface{} {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	var result []map[string]interface{}
	for _, row := range store.data {
		result = append(result, row)
	}
	return result
}

func (store *KeyValueStore) ClearCache() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.cache = make(map[interface{}]map[string]interface{})
}

type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %v not found", e.Key)
}
