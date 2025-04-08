package store

import (
	"fmt"
	"sync"
	"solune/data"
)

type KeyValueStore struct {
	data  map[interface{}][]byte
	cache map[interface{}][]byte
	mutex sync.RWMutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data:  make(map[interface{}][]byte),
		cache: make(map[interface{}][]byte),
	}
}

func (store *KeyValueStore) Set(key interface{}, value map[string]interface{}) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	binValue, err := data.MapToBinary(value)
	if err != nil {
		return err
	}

	store.cache[key] = binValue
	store.data[key] = binValue
	return nil
}

func (store *KeyValueStore) Get(key interface{}) (map[string]interface{}, error) {
	if binValue, found := store.cache[key]; found {
		return data.BinaryToMap(binValue)
	}

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if binValue, exists := store.data[key]; exists {
		store.cache[key] = binValue
		return data.BinaryToMap(binValue)
	}

	return nil, &KeyNotFoundError{Key: key}
}


func (store *KeyValueStore) GetAllData() []map[string]interface{} {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	var result []map[string]interface{}

	for key, binValue := range store.data {
		if m, err := data.BinaryToMap(binValue); err == nil {
			m["key"] = key
			result = append(result, m)
		}
	}
	return result
}


func (store *KeyValueStore) ClearCache() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.cache = make(map[interface{}][]byte)
}

type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %v not found", e.Key)
}
