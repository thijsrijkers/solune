package store

import (
	"encoding/base64"
	"fmt"
	"solune/filestore"
	"sort"
	"sync"
)

const shards = 50

type Shard struct {
	data map[int][]byte
	mu   sync.RWMutex
}

type KeyValueStore struct {
	shards    [shards]Shard
	fileStore *filestore.FileStore
}

func NewKeyValueStore(fs *filestore.FileStore) *KeyValueStore {
	store := &KeyValueStore{fileStore: fs}
	for i := range store.shards {
		store.shards[i].data = make(map[int][]byte)
	}
	return store
}

func (store *KeyValueStore) getShard(key int) *Shard {
	return &store.shards[key%shards]
}

func (store *KeyValueStore) Set(key int, value string) error {
	shard := store.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = []byte(value)

	go func() {
		encoded := base64.StdEncoding.EncodeToString([]byte(value))
		if err := store.fileStore.Update(fmt.Sprintf("%d", key), encoded); err != nil {
			fmt.Printf("[ERROR] Failed to write key %d to filestore: %v\n", key, err)
		}
	}()

	return nil
}

func (store *KeyValueStore) Update(key int, newValue string) error {
	return store.Set(key, newValue)
}

func (store *KeyValueStore) Get(key int) (string, error) {
	shard := store.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	val, ok := shard.data[key]
	if !ok {
		return "", &KeyNotFoundError{Key: key}
	}
	return string(val), nil
}

func (store *KeyValueStore) Delete(key int) error {
	shard := store.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.data[key]; !ok {
		return &KeyNotFoundError{Key: key}
	}
	delete(shard.data, key)

	if err := store.fileStore.Delete(fmt.Sprintf("%d", key)); err != nil {
		return err
	}
	return nil
}

func (store *KeyValueStore) GetAllData() map[int]string {
	result := make(map[int]string)

	keys := make([]int, 0)
	for i := range store.shards {
		shard := &store.shards[i]
		shard.mu.RLock()
		for k, v := range shard.data {
			keys = append(keys, k)
			result[k] = string(v)
		}
		shard.mu.RUnlock()
	}

	sort.Ints(keys)

	sorted := make(map[int]string, len(keys))
	for _, k := range keys {
		sorted[k] = result[k]
	}
	return sorted
}

type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %v not found", e.Key)
}
