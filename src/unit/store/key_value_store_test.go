package store_test

import (
	"errors"
	"reflect"
	"testing"
	"paper/src/store"
)

func TestKeyValueStore(test *testing.T) {
	validate := func(value interface{}) error {
		if str, ok := value.(string); ok && len(str) < 3 {
			return errors.New("string length must be at least 3 characters")
		}
		return nil
	}

	store := store.NewKeyValueStore(reflect.TypeOf(""), reflect.TypeOf(""), validate)
	if err := store.Set("key1", "go"); err == nil {
		test.Errorf("Expected validation error, got nil")
	}

	if err := store.Set("key1", "paper"); err != nil {
		test.Errorf("Unexpected error: %v", err)
	}

	value, err := store.Get("key1")
	if err != nil || value != "paper" {
		test.Errorf("Expected 'paper', got '%v' with error '%v'", value, err)
	}
}
