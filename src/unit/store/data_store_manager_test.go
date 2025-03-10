package store_test

import (
	"reflect"
	"testing"
	"paper/src/store"
)

func TestKeyValueStoreManager(test *testing.T) {
	manager := store.NewDataStoreManager()

	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	store := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)
	manager.AddStore("testStore", store)

	retrievedStore, exists := manager.GetStore("testStore")
	if !exists || retrievedStore != store {
		test.Errorf("Expected to retrieve the correct store, but got different result")
	}
}
