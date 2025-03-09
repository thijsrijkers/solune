package store

import (
	"reflect"
)

type Schema struct {
	KeyType   reflect.Type
	ValueType reflect.Type
	Validate  func(interface{}) error
}
