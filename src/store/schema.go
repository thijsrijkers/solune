package store

import (
	"errors"
	"fmt"
	"reflect"
)

type ColumnSchema map[string]reflect.Type

type Schema struct {
	KeyType     reflect.Type
	ColumnTypes ColumnSchema
	Validate    func(map[string]interface{}) error
}

func (schema *Schema) ValidateRow(value map[string]interface{}) error {
	for column, value := range value {
		expectedType, exists := schema.ColumnTypes[column]
		if !exists {
			return fmt.Errorf("invalid column: %s", column)
		}
		if reflect.TypeOf(value) != expectedType {
			return fmt.Errorf("invalid type for column %s: expected %v, got %v", column, expectedType, reflect.TypeOf(value))
		}
	}
	if schema.Validate != nil {
		if err := schema.Validate(value); err != nil {
			return fmt.Errorf("custom validation failed: %w", err)
		}
	}
	return nil
}

func (schema *Schema) ValidateKey(key interface{}) error {
	if reflect.TypeOf(key) != schema.KeyType {
		return errors.New("invalid key type")
	}
	return nil
}
