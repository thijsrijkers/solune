package store

import (
	"fmt"
	"sync"
	"encoding/base64"
	"solune/data"
    "solune/filestore"
)

type KeyValueStore struct {
	data  map[interface{}][]byte
	cache map[interface{}][]byte
	mutex sync.RWMutex
	fileStore *filestore.FileStore
}

func NewKeyValueStore(fs *filestore.FileStore) *KeyValueStore {
	return &KeyValueStore{
		data:  make(map[interface{}][]byte),
		cache: make(map[interface{}][]byte),
		fileStore: fs,
	}
}

func (store *KeyValueStore) Set(key interface{}, value map[string]interface{}) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	binValue, err := data.MapToBinary(value)
	if err != nil {
		return err
	}

	
	if err := store.fileStore.Update(fmt.Sprintf("%v", key), base64.StdEncoding.EncodeToString(binValue)); err != nil {
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

func (store *KeyValueStore) Update(key interface{}, newValue map[string]interface{}) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, exists := store.data[key]
	if !exists {
		return &KeyNotFoundError{Key: key}
	}

	binValue, err := data.MapToBinary(newValue)
	if err != nil {
		return err
	}

	if err := store.fileStore.Update(fmt.Sprintf("%v", key), base64.StdEncoding.EncodeToString(binValue)); err != nil {
		return err
	}

	store.data[key] = binValue
	store.cache[key] = binValue
	return nil
}

func (store *KeyValueStore) Delete(key interface{}) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, exists := store.data[key]; !exists {
		return &KeyNotFoundError{Key: key}
	}

	if err := store.fileStore.Delete(fmt.Sprintf("%v", key)); err != nil {
		return err
	}

	delete(store.data, key)
	delete(store.cache, key)
	return nil
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
