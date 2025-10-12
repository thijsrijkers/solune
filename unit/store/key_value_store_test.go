package store_test

import (
	"testing"

	"solune/store"
	"solune/filestore"
)

func TestKeyValueStore(t *testing.T) {
	fs, err := filestore.New("testKeyValueStore", "9000")
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	kv := store.NewKeyValueStore(fs)

	key1 := 1
	value1 := `{"name":"Alice","age":30}`

	t.Run("Set and Get", func(t *testing.T) {
		err := kv.Set(key1, value1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		got, err := kv.Get(key1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if got != value1 {
			t.Errorf("expected value %v, got %v", value1, got)
		}
	})

	t.Run("GetAllData", func(t *testing.T) {
		key2 := 2
		value2 := `{"name":"Bob","age":25}`

		err := kv.Set(key2, value2)
		if err != nil {
			t.Errorf("expected no error on Set, got %v", err)
		}

		allData := kv.GetAllData()

		if len(allData) != 2 {
			t.Errorf("expected 2 records, got %d", len(allData))
		}

		if allData[key1] != value1 || allData[key2] != value2 {
			t.Errorf("expected stored values to match, got %v", allData)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newValue := `{"name":"Alice Updated","age":31}`

		err := kv.Update(key1, newValue)
		if err != nil {
			t.Errorf("expected no error on update, got %v", err)
		}

		got, err := kv.Get(key1)
		if err != nil {
			t.Errorf("expected no error on get after update, got %v", err)
		}

		if got != newValue {
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
