package store

import (
	"fmt"
	"sync"
	"encoding/base64"
	"github.com/google/uuid"
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

    var keyStr string
    switch k := key.(type) {
    case string:
        keyStr = k
    case fmt.Stringer:
        keyStr = k.String()
    }

    value["key"] = keyStr

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
	var validKey string

	switch v := key.(type) {
	case string:
		validKey = v
	case uuid.UUID:
		validKey = v.String()
	default:
		return nil, fmt.Errorf("invalid key type: expected string or uuid.UUID, got %T", key)
	}

	if binValue, found := store.cache[validKey]; found {
		return data.BinaryToMap(binValue)
	}

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	if binValue, exists := store.data[validKey]; exists {
		store.cache[validKey] = binValue 
		return data.BinaryToMap(binValue)
	}

	return nil, &KeyNotFoundError{Key: validKey}
}

func (store *KeyValueStore) Update(key interface{}, newValue map[string]interface{}) error {
	var validKey string

	switch v := key.(type) {
	case string:
		validKey = v
	case uuid.UUID:
		validKey = v.String()
	default:
		return fmt.Errorf("invalid key type: expected string or uuid.UUID, got %T", key)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, exists := store.data[validKey]
	if !exists {
		return &KeyNotFoundError{Key: validKey}
	}

	binValue, err := data.MapToBinary(newValue)
	if err != nil {
		return err
	}

	if err := store.fileStore.Update(fmt.Sprintf("%v", validKey), base64.StdEncoding.EncodeToString(binValue)); err != nil {
		return err
	}

	store.data[validKey] = binValue
	store.cache[validKey] = binValue
	return nil
}

func (store *KeyValueStore) Delete(key interface{}) error {
	var validKey string

	switch v := key.(type) {
	case string:
		validKey = v
	case uuid.UUID:
		validKey = v.String()
	default:
		return fmt.Errorf("invalid key type: expected string or uuid.UUID, got %T", key)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, exists := store.data[validKey]; !exists {
		return &KeyNotFoundError{Key: validKey}
	}

	if err := store.fileStore.Delete(fmt.Sprintf("%v", validKey)); err != nil {
		return err
	}

	delete(store.data, validKey)
	delete(store.cache, validKey)
	return nil
}

func (store *KeyValueStore) GetAllData() []map[string]interface{} {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	result := make([]map[string]interface{}, 0, len(store.data))

	for _, binValue := range store.data {
		m, err := data.BinaryToMap(binValue)
		if err != nil {
			continue
		}
		result = append(result, m)
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
