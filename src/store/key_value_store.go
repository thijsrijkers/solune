package store

import (
	"errors"
	"fmt"
	"reflect"
)

type KeyValueStore struct {
	data   map[interface{}]interface{}
	schema Schema
}

func NewKeyValueStore(keyType, valueType reflect.Type, validate func(interface{}) error) *KeyValueStore {
	return &KeyValueStore{
		data: make(map[interface{}]interface{}),
		schema: Schema{
			KeyType:   keyType,
			ValueType: valueType,
			Validate:  validate,
		},
	}
}

func (store *KeyValueStore) Set(key, value interface{}) error {
	if reflect.TypeOf(key) != store.schema.KeyType {
		return errors.New("invalid key type")
	}
	if reflect.TypeOf(value) != store.schema.ValueType {
		return errors.New("invalid value type")
	}
	if store.schema.Validate != nil {
		if err := store.schema.Validate(value); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}
	store.data[key] = value
	return nil
}

func (store *KeyValueStore) Get(key interface{}) (interface{}, error) {
	if reflect.TypeOf(key) != store.schema.KeyType {
		return nil, errors.New("invalid key type")
	}
	value, exists := store.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}
	return value, nil
}
