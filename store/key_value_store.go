package store

import (
	"encoding/base64"
	"fmt"
	"solune/filestore"
	"sort"
	"sync"
	"sync/atomic"
)

type Shard struct {
	data map[int][]byte
	mu   sync.RWMutex
}

type KeyValueStore struct {
	shards     []Shard
	shardCount int
	fileStore  *filestore.FileStore
	NextKey    atomic.Int64
}

func NewKeyValueStore(fs *filestore.FileStore, numberOfShards ...int) *KeyValueStore {
	shardCount := 50
	if len(numberOfShards) > 0 && numberOfShards[0] > 0 {
		shardCount = numberOfShards[0]
	}
	store := &KeyValueStore{
		fileStore:  fs,
		shards:     make([]Shard, shardCount),
		shardCount: shardCount,
	}

	store.NextKey.Store(1)

	for i := range store.shards {
		store.shards[i].data = make(map[int][]byte)
	}

	return store
}

func (store *KeyValueStore) getShard(key int) *Shard {
	return &store.shards[key%store.shardCount]
}

func (store *KeyValueStore) Set(key int, value string) error {
	shard := store.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = []byte(value)

	// NextKey always tracks the highest key seen + 1.
	// A Compare And Swap (CAS) loop is used instead of a simple Store because multiple goroutines
	// across different shards can call Set concurrently. If two goroutines both
	// load the same current value and race to update it, only one CAS wins,
	// the other retries with the latest value until it either succeeds or finds
	// that NextKey is already higher than its key.
	for {
		current := store.NextKey.Load()
		if int64(key) < current {
			break
		}
		if store.NextKey.CompareAndSwap(current, int64(key)+1) {
			break
		}
	}

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
