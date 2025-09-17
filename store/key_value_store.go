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
	data sync.Map 
	fileStore *filestore.FileStore
}

func NewKeyValueStore(fs *filestore.FileStore) *KeyValueStore {
	return &KeyValueStore{
		data: sync.Map{},
		fileStore: fs,
	}
}

func (store *KeyValueStore) Set(key interface{}, value map[string]interface{}) error {
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
	
	if err := store.fileStore.Update(fmt.Sprintf("%v", keyStr), base64.StdEncoding.EncodeToString(binValue)); err != nil {
		return err
	}

    store.data.Store(keyStr, binValue)
	return nil
}

func (store *KeyValueStore) Get(key string) (map[string]interface{}, error) {
    if binValue, ok := store.data.Load(key); ok {
        return data.BinaryToMap(binValue.([]byte))
    }
    return nil, &KeyNotFoundError{Key: key}
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

	binValue, err := data.MapToBinary(newValue)
	if err != nil {
		return err
	}

	if err := store.fileStore.Update(fmt.Sprintf("%v", validKey), base64.StdEncoding.EncodeToString(binValue)); err != nil {
		return err
	}

    store.data.Store(key, binValue)
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

	if _, exists := store.data.Load(validKey); !exists {
		return &KeyNotFoundError{Key: validKey}
	}

	if err := store.fileStore.Delete(fmt.Sprintf("%v", validKey)); err != nil {
		return err
	}

	store.data.Delete(key)
	return nil
}

func (store *KeyValueStore) GetAllData() []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 128)

	store.data.Range(func(_, value interface{}) bool {
		if binValue, ok := value.([]byte); ok {
			if m, err := data.BinaryToMap(binValue); err == nil {
				result = append(result, m)
			}
		}
		return true
	})

	return result
}

type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %v not found", e.Key)
}
