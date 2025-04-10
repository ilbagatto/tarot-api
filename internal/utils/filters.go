package utils

import (
	"fmt"
	"sort"
	"strings"
)

// BuildWhereClause constructs a SQL WHERE clause from a map of filters.
// Example:
// filters := map[string]any{"number": 0, "position": "straight"}
// returns: "WHERE number = $1 AND position = $2", []any{0, "straight"}
func BuildWhereClause(filters map[string]any, startIndex int) (string, []any) {
	var clauses []string
	var args []any
	index := startIndex

	// Стабилизируем порядок обхода
	keys := make([]string, 0, len(filters))
	for key := range filters {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Обходим ключи по порядку
	for _, key := range keys {
		value := filters[key]
		if value != nil {
			clauses = append(clauses, fmt.Sprintf("%s = $%d", key, index))
			args = append(args, value)
			index++
		}
	}

	if len(clauses) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}
