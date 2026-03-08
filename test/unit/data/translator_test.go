package data_test

import (
	"testing"
	"reflect"
	"solune/data"
)

func TestMapToBinary(t *testing.T) {
	originalMap := map[string]interface{}{
		"Name":   "Alice",
		"Age":    30,
		"Active": true,
	}

	binaryData, err := data.MapToBinary(originalMap)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(binaryData) == 0 {
		t.Fatal("Expected binary data to be non-empty")
	}
}

func TestBinaryToMap(t *testing.T) {
	originalMap := map[string]interface{}{
		"Name":   "Alice",
		"Age":    30,
		"Active": true,
	}

	binaryData, err := data.MapToBinary(originalMap)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decodedMap, err := data.BinaryToMap(binaryData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(originalMap, decodedMap) {
		t.Errorf("Expected map %v, but got %v", originalMap, decodedMap)
	}
}

func TestBinaryToMapInvalidData(t *testing.T) {
	_, err := data.BinaryToMap([]byte{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
