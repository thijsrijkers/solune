package sql

import (
	"fmt"
	"paper/src/store"
	"strings"
)

type SQLDataTranslator struct {
	manager *store.DataStoreManager
}

func NewSQLDataTranslator(manager *store.DataStoreManager) *SQLDataTranslator {
	return &SQLDataTranslator{manager: manager}
}

// Translate processes the SQL query and retrieves the corresponding rows
func (translator *SQLDataTranslator) Translate(query string) ([]map[string]interface{}, error) {
	// Use the parseSQL function to break down the query
	columns, tableName, err := parseSQL(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL query: %v", err)
	}

	// Normalize the table name to lowercase to avoid case-sensitivity issues
	tableName = strings.ToLower(tableName)

	// Retrieve the KeyValueStore for the specified table
	store, exists := translator.manager.GetStore(tableName)
	if !exists {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}

	// Retrieve rows from the KeyValueStore
	var result []map[string]interface{}
	for _, row := range store.GetAllData() {
		// Select only the columns requested in the query
		selectedRow := make(map[string]interface{})
		for _, column := range columns {
			if value, ok := row[column]; ok {
				selectedRow[column] = value
			}
		}
		result = append(result, selectedRow)
	}

	return result, nil
}
