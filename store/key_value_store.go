package store

import (
	"encoding/base64"
	"fmt"
	"solune/filestore"
	"sync"

	"github.com/google/btree"
)

type item struct {
	key   int
	value []byte
}

func (a item) Less(b btree.Item) bool {
	return a.key < b.(*item).key
}

type KeyValueStore struct {
	tree      *btree.BTree
	fileStore *filestore.FileStore
	mu        sync.RWMutex
}

func NewKeyValueStore(fs *filestore.FileStore) *KeyValueStore {
	return &KeyValueStore{
		tree:      btree.New(2),
		fileStore: fs,
	}
}

func (store *KeyValueStore) Set(key int, value string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	binValue := []byte(value)
	it := &item{key: key, value: binValue}
	store.tree.ReplaceOrInsert(it)

	go func(key int, binValue []byte) {
		encoded := base64.StdEncoding.EncodeToString(binValue)
		if err := store.fileStore.Update(fmt.Sprintf("%d", key), encoded); err != nil {
			fmt.Printf("[ERROR] Failed to write key %d to filestore: %v\n", key, err)
		}
	}(key, binValue)

	return nil
}

func (store *KeyValueStore) Update(key int, newValue string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	binValue := []byte(newValue)
	it := &item{key: key, value: binValue}
	store.tree.ReplaceOrInsert(it)

	go func(key int, binValue []byte) {
		encoded := base64.StdEncoding.EncodeToString(binValue)
		if err := store.fileStore.Update(fmt.Sprintf("%d", key), encoded); err != nil {
			fmt.Printf("[ERROR] Failed to update key %d in filestore: %v\n", key, err)
		}
	}(key, binValue)

	return nil
}

func (store *KeyValueStore) Get(key int) (string, error) {
	it := store.tree.Get(&item{key: key})
	if it == nil {
		return "", &KeyNotFoundError{Key: key}
	}

	return string(it.(*item).value), nil
}

func (store *KeyValueStore) Delete(key int) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	it := store.tree.Delete(&item{key: key})
	if it == nil {
		return &KeyNotFoundError{Key: key}
	}

	if err := store.fileStore.Delete(fmt.Sprintf("%d", key)); err != nil {
		return err
	}

	return nil
}

func (store *KeyValueStore) GetAllData() map[int]string {
	store.mu.RLock()
	defer store.mu.RUnlock()

	result := make(map[int]string)
	store.tree.Ascend(func(i btree.Item) bool {
		it := i.(*item)
		result[it.key] = string(it.value)
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
