package store

import (
	"reflect"
	"errors"
)

type KeyValueStore struct {
	data   map[interface{}]map[string]interface{}
	schema Schema
}

func NewKeyValueStore(keyType reflect.Type, columnTypes ColumnSchema, validate func(map[string]interface{}) error) *KeyValueStore {
	return &KeyValueStore{
		data: make(map[interface{}]map[string]interface{}),
		schema: Schema{
			KeyType:     keyType,
			ColumnTypes: columnTypes,
			Validate:    validate,
		},
	}
}

func (store *KeyValueStore) Set(key interface{}, value map[string]interface{}) error {
	if err := store.schema.ValidateKey(key); err != nil {
		return err
	}
	if err := store.schema.ValidateRow(value); err != nil {
		return err
	}
	store.data[key] = value
	return nil
}

func (store *KeyValueStore) Get(key interface{}) (map[string]interface{}, error) {
	if err := store.schema.ValidateKey(key); err != nil {
		return nil, err
	}
	value, exists := store.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (store *KeyValueStore) GetAllData() []map[string]interface{} {
	var result []map[string]interface{}
	for _, row := range store.data {
		result = append(result, row)
	}
	return result
}
