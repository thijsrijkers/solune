package sql_test

import (
	"testing"
	"stack/src/sql"
	"stack/src/store"
	"reflect"
)

func setupDataStore() *store.DataStoreManager {
	manager := store.NewDataStoreManager()

	// Define column types for the `users` table (assuming ID and Name are the columns)
	columnTypes := store.ColumnSchema{
		"ID":   reflect.TypeOf(int(0)),
		"Name": reflect.TypeOf(""),
	}

	// Create a KeyValueStore for the `users` table
	userStore := store.NewKeyValueStore(reflect.TypeOf(0), columnTypes, nil)

	// Add some test data to the `users` table
	userStore.Set(1, map[string]interface{}{"ID": 1, "Name": "John"})
	userStore.Set(2, map[string]interface{}{"ID": 2, "Name": "Alice"})
	userStore.Set(3, map[string]interface{}{"ID": 3, "Name": "Bob"})

	// Add the `users` store to the manager
	manager.AddStore("users", userStore)

	return manager
}

func TestSQLTranslator_Select(t *testing.T) {
	// Set up the DataStoreManager and add the users table
	manager := setupDataStore()

	// Create the SQLTranslator
	translator := sql.NewSQLDataTranslator(manager)

	tests := []struct {
		query      string
		expected   int // Expected number of rows returned
		expectFail bool // Whether we expect the query to fail
	}{
		{
			query:    "SELECT ID, Name FROM users",
			expected: 3, // We expect 3 rows (John, Alice, Bob)
		},
		{
			query:    "SELECT Name FROM users",
			expected: 3, // We expect 3 rows (John, Alice, Bob)
		},
		{
			query:    "SELECT ID FROM users",
			expected: 3, // We expect 3 rows (John, Alice, Bob)
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			rows, err := translator.Translate(tt.query)
			if err != nil {
				if !tt.expectFail {
					t.Errorf("Translate(%s) failed: %v", tt.query, err)
				}
				return
			}

			if len(rows) != tt.expected {
				t.Errorf("Expected %d rows, got %d", tt.expected, len(rows))
			}
		})
	}
}