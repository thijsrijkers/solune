package store_test

import (
	"errors"
	"reflect"
	"testing"
	"paper/src/store"
)

func TestNewKeyValueStore(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)

	if db == nil {
		t.Fatalf("expected non-nil KeyValueStore instance")
	}
}

func TestSetAndGet(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)

	err := db.Set("user1", map[string]interface{}{"age": 30, "name": "Alice"})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	value, err := db.Get("user1")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if value["name"] != "Alice" || value["age"].(int) != 30 {
		t.Errorf("expected {name: Alice, age: 30}, got: %v", value)
	}
}

func TestSetInvalidKey(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)

	err := db.Set(123, map[string]interface{}{"age": 30, "name": "Alice"})
	if err == nil {
		t.Errorf("expected error for invalid key type, got nil")
	}
}

func TestSetInvalidValueType(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)

	err := db.Set("user1", map[string]interface{}{"age": "thirty", "name": "Alice"})
	if err == nil {
		t.Errorf("expected error for invalid value type, got nil")
	}
}

func TestSetWithCustomValidation(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	validate := func(row map[string]interface{}) error {
		if name, ok := row["name"].(string); ok && name == "" {
			return errors.New("name cannot be empty")
		}
		return nil
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, validate)

	err := db.Set("user1", map[string]interface{}{"age": 30, "name": ""})
	if err == nil {
		t.Errorf("expected error for empty name, got nil")
	}
}

func TestGetNonExistentKey(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	db := store.NewKeyValueStore(reflect.TypeOf(""), columnTypes, nil)

	_, err := db.Get("nonexistent")
	if err == nil {
		t.Errorf("expected error for missing key, got nil")
	}
}
