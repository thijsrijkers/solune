package store

import (
	"encoding/base64"
	"fmt"
	"solune/filestore"
	"sort"
	"strconv"
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

	// Seed NextKey from the highest key already on disk.
	maxKey := int64(0)
	allKeys := fs.Keys() // see below
	for _, k := range allKeys {
		if parsed, err := strconv.ParseInt(k, 10, 64); err == nil {
			if parsed > maxKey {
				maxKey = parsed
			}
		}
	}
	store.NextKey.Store(maxKey + 1)

	return store
}

func (store *KeyValueStore) getShard(key int) *Shard {
	idx := key % store.shardCount
	if idx < 0 {
		idx += store.shardCount
	}
	return &store.shards[idx]
}

func (store *KeyValueStore) Set(key int, value string) error {
	shard := store.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	if err := store.fileStore.Update(fmt.Sprintf("%d", key), encoded); err != nil {
		fmt.Printf("[ERROR] Failed to write key %d to filestore: %v\n", key, err)
		return err
	}

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

	return nil
}

func (store *KeyValueStore) Update(key int, newValue string) error {
	return store.Set(key, newValue)
}

func (store *KeyValueStore) Get(key int) (string, error) {
	shard := store.getShard(key)
	shard.mu.RLock()
	val, ok := shard.data[key]
	shard.mu.RUnlock()
	if ok {
		return string(val), nil
	}

	encoded, err := store.fileStore.Get(strconv.Itoa(key))

	if err != nil {
		return "", &KeyNotFoundError{Key: key}
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode value for key %d: %w", key, err)
	}

	shard.mu.Lock()
	shard.data[key] = decoded
	shard.mu.Unlock()

	return string(decoded), nil
}

func (store *KeyValueStore) Delete(key int) error {
	shard := store.getShard(key)
	shard.mu.Lock()
	_, inMemory := shard.data[key]
	if inMemory {
		delete(shard.data, key)
	}
	shard.mu.Unlock()

	if err := store.fileStore.Delete(fmt.Sprintf("%d", key)); err != nil {
		if !inMemory {
			return &KeyNotFoundError{Key: key}
		}
		return err
	}
	return nil
}

func (store *KeyValueStore) GetAllData() map[int]string {
	result := make(map[int]string)
	keys := make([]int, 0)
	seen := make(map[int]struct{})

	allFromDisk, err := store.fileStore.All()
	if err == nil {
		for keyStr, encoded := range allFromDisk {
			k, parseErr := strconv.Atoi(keyStr)
			if parseErr != nil {
				continue
			}
			decoded, decodeErr := base64.StdEncoding.DecodeString(encoded)
			if decodeErr != nil {
				continue
			}
			result[k] = string(decoded)
			seen[k] = struct{}{}
			keys = append(keys, k)
		}
	}

	for i := range store.shards {
		shard := &store.shards[i]
		shard.mu.RLock()
		for k, v := range shard.data {
			if _, exists := seen[k]; !exists {
				keys = append(keys, k)
				seen[k] = struct{}{}
			}
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
