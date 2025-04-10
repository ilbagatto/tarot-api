package utils_test

import (
	"testing"

	"github.com/ilbagatto/tarot-api/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseBoolParam(t *testing.T) {
	assert.True(t, utils.ParseBoolParam("true"))
	assert.True(t, utils.ParseBoolParam("1"))
	assert.True(t, utils.ParseBoolParam("YES"))
	assert.False(t, utils.ParseBoolParam("false"))
	assert.False(t, utils.ParseBoolParam("0"))
	assert.False(t, utils.ParseBoolParam(""))
}
