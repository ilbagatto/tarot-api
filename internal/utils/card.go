package utils

import (
	"os"
	"strings"
)

// GetCardImage constructs a full URL to a card image.
// It returns nil if BASE_URL is not set.
func GetCardImage(path string) *string {
	base := os.Getenv("BASE_URL")
	if base == "" {
		return nil
	}

	var url string
	if strings.HasPrefix(path, "/") {
		url = base + path
	} else {
		url = base + "/" + path
	}

	return &url
}
