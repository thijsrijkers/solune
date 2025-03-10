package sql

import (
	"fmt"
	"regexp"
	"strings"
)

func parseSQL(query string) ([]string, string, error) {
	query = strings.ToUpper(strings.TrimSpace(query))

	selectPattern := `SELECT\s+(.*)\s+FROM\s+([a-zA-Z_][a-zA-Z0-9_]*)`
	re, err := regexp.Compile(selectPattern)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compile regex: %v", err)
	}

	matches := re.FindStringSubmatch(query)
	if len(matches) < 3 {
		return nil, "", fmt.Errorf("query format invalid, expected SELECT ... FROM ...")
	}

	columnsPart := matches[1]
	tableName := matches[2]

	columns := strings.Split(columnsPart, ",")
	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	return columns, tableName, nil
}