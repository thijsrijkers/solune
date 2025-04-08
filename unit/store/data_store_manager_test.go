package store_test

import (
	"reflect"
	"testing"
	"solune/store"
)

func TestDataStoreManager(t *testing.T) {
	manager := store.NewDataStoreManager()

	manager.AddStore("users")
	store, exists := manager.GetStore("users")
	if !exists {
		t.Fatalf("expected store 'users' to exist")
	}

	key := "user1"
	value := map[string]interface{}{"name": "Alice", "age": 30}
	err := store.Set(key, value)
	if err != nil {
		t.Errorf("expected no error from Set, got %v", err)
	}

	got, err := store.Get(key)
	if err != nil {
		t.Errorf("expected no error from Get, got %v", err)
	}
	if !reflect.DeepEqual(got, value) {
		t.Errorf("expected value %v, got %v", value, got)
	}

	all := store.GetAllData()
	if len(all) != 1 {
		t.Errorf("expected 1 record in GetAllData, got %d", len(all))
	}
	if !reflect.DeepEqual(all[0], value) {
		t.Errorf("expected GetAllData to return %v, got %v", value, all[0])
	}

	manager.AddStore("products")
	_, exists = manager.GetStore("products")
	if !exists {
		t.Errorf("expected store 'products' to exist")
	}
}
