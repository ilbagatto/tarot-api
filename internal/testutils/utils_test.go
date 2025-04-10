package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString_LengthRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandomString(10, 30)
		assert.GreaterOrEqual(t, len(s), 10)
		assert.LessOrEqual(t, len(s), 30)
	}
}

func TestRandomString_OnlyCharset(t *testing.T) {
	allowed := map[rune]bool{}
	for _, ch := range charset {
		allowed[ch] = true
	}

	s := RandomString(20, 20)
	for _, ch := range s {
		assert.True(t, allowed[ch], "invalid character: %q", ch)
	}
}

func TestRandomString_NotEmpty(t *testing.T) {
	s := RandomString(10, 10)
	assert.NotEmpty(t, s)
}

func TestRandomString_Repeatable(t *testing.T) {
	s1 := RandomString(10, 10)
	s2 := RandomString(10, 10)
	assert.NotEqual(t, s1, s2, "consecutive strings should likely differ")
}

func TestRandomString_MinEqualsMax(t *testing.T) {
	s := RandomString(15, 15)
	assert.Equal(t, 15, len(s))
}
