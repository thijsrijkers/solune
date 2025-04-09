package store_test

import (
	"reflect"
	"testing"
	"solune/store"
)

func TestDataStoreManager(t *testing.T) {
	manager := store.NewDataStoreManager()

	manager.AddStore("users")
	userStore, exists := manager.GetStore("users")
	if !exists {
		t.Fatalf("expected store 'users' to exist")
	}

	key := "user1"
	value := map[string]interface{}{"name": "Alice", "age": 30}
	err := userStore.Set(key, value)
	if err != nil {
		t.Errorf("expected no error from Set, got %v", err)
	}

	got, err := userStore.Get(key)
	if err != nil {
		t.Errorf("expected no error from Get, got %v", err)
	}
	if !reflect.DeepEqual(normalize(got), normalize(value)) {
		t.Errorf("expected value %v, got %v", value, got)
	}

	all := userStore.GetAllData()
	if len(all) != 1 {
		t.Errorf("expected 1 record in GetAllData, got %d", len(all))
	}
	expected := normalize(map[string]interface{}{"key": "user1", "name": "Alice", "age": 30})
	if !reflect.DeepEqual(normalize(all[0]), expected) {
		t.Errorf("expected GetAllData to return %v, got %v", expected, all[0])
	}

	manager.AddStore("products")
	_, exists = manager.GetStore("products")
	if !exists {
		t.Errorf("expected store 'products' to exist")
	}

	removed := manager.RemoveStore("users")
	if !removed {
		t.Errorf("expected RemoveStore to return true for existing store")
	}

	_, exists = manager.GetStore("users")
	if exists {
		t.Errorf("expected 'users' store to be removed")
	}

	removed = manager.RemoveStore("nonexistent")
	if removed {
		t.Errorf("expected RemoveStore to return false for non-existent store")
	}
}
