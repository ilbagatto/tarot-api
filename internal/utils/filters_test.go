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

	require.Equal(t, "WHERE source = $1 AND position = $2", clause)
	require.Equal(t, []any{42, "straight"}, args)
}
