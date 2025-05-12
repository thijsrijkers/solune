package store_test

import (
	"reflect"
	"testing"
	"solune/store"
	"solune/filestore"
)

func normalize(m map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})
	for k, v := range m {
		switch v := v.(type) {
		case int:
			normalized[k] = float64(v)
		case int64:
			normalized[k] = float64(v)
		case float32:
			normalized[k] = float64(v)
		default:
			normalized[k] = v
		}
	}
	return normalized
}

func TestKeyValueStore(t *testing.T) {
	// Initialize the file store with a test filename
	fs, err := filestore.New("testKeyValueStore")
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	kv := store.NewKeyValueStore(fs)

	key1 := "user1"
	value1 := map[string]interface{}{"name": "Alice", "age": 30}

	t.Run("Set and Get", func(t *testing.T) {
		err := kv.Set(key1, value1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		got, err := kv.Get(key1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(normalize(got), normalize(value1)) {
			t.Errorf("expected value %v, got %v", value1, got)
		}
	})

	t.Run("GetAllData", func(t *testing.T) {
		key2 := "user2"
		value2 := map[string]interface{}{"name": "Bob", "age": 25}
		kv.Set(key2, value2)

		got := kv.GetAllData()

		if len(got) != 2 {
			t.Errorf("expected 2 records, got %d", len(got))
		}

		expected1 := normalize(map[string]interface{}{"key": key1, "name": "Alice", "age": 30})
		expected2 := normalize(map[string]interface{}{"key": key2, "name": "Bob", "age": 25})

		found1, found2 := false, false
		for _, row := range got {
			if reflect.DeepEqual(normalize(row), expected1) {
				found1 = true
			}
			if reflect.DeepEqual(normalize(row), expected2) {
				found2 = true
			}
		}
		if !found1 || !found2 {
			t.Errorf("expected both values in the result")
		}
	})

	t.Run("ClearCache", func(t *testing.T) {
		key3 := "user3"
		value3 := map[string]interface{}{"name": "Charlie", "age": 35}
		kv.Set(key3, value3)

		kv.ClearCache()

		got, err := kv.Get(key3)
		if err != nil {
			t.Errorf("expected no error after cache clear, got %v", err)
		}
		if !reflect.DeepEqual(normalize(got), normalize(value3)) {
			t.Errorf("expected value %v, got %v", value3, got)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newValue := map[string]interface{}{"name": "Alice Updated", "age": 31}
		err := kv.Update(key1, newValue)
		if err != nil {
			t.Errorf("expected no error on update, got %v", err)
		}

		got, err := kv.Get(key1)
		if err != nil {
			t.Errorf("expected no error on get after update, got %v", err)
		}

		if !reflect.DeepEqual(normalize(got), normalize(newValue)) {
			t.Errorf("expected updated value %v, got %v", newValue, got)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := kv.Delete(key1)
		if err != nil {
			t.Errorf("expected no error on delete, got %v", err)
		}

		_, err = kv.Get(key1)
		if err == nil {
			t.Errorf("expected error when getting deleted key, got nil")
		}
	})
}
