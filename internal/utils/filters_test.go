package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildWhereClause(t *testing.T) {
	clause, args := BuildWhereClause(map[string]any{
		"source":   42,
		"position": "straight",
	}, 1)

	require.Equal(t, "WHERE position = $1 AND source = $2", clause)
	require.Equal(t, []any{"straight", 42}, args)
}
