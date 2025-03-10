package store_test

import (
	"errors"
	"reflect"
	"testing"
	"paper/src/store"
)

func TestValidateKey(t *testing.T) {
	schema := store.Schema{KeyType: reflect.TypeOf("")}

	err := schema.ValidateKey("validKey")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	err = schema.ValidateKey(123)
	if err == nil {
		t.Errorf("expected error for invalid key type, got nil")
	}
}

func TestValidateRow(t *testing.T) {
	columnTypes := store.ColumnSchema{
		"age":  reflect.TypeOf(0),
		"name": reflect.TypeOf(""),
	}

	schema := store.Schema{
		ColumnTypes: columnTypes,
		Validate: func(row map[string]interface{}) error {
			if name, ok := row["name"].(string); ok && name == "" {
				return errors.New("name cannot be empty")
			}
			return nil
		},
	}

	validRow := map[string]interface{}{"age": 25, "name": "John"}
	err := schema.ValidateRow(validRow)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	invalidRow := map[string]interface{}{"age": "twenty-five", "name": "John"}
	err = schema.ValidateRow(invalidRow)
	if err == nil {
		t.Errorf("expected error for invalid column type, got nil")
	}

	emptyNameRow := map[string]interface{}{"age": 30, "name": ""}
	err = schema.ValidateRow(emptyNameRow)
	if err == nil {
		t.Errorf("expected error for empty name, got nil")
	}
}