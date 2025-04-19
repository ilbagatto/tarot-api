package utils

import (
	"os"
	"strings"
)

// GetImageURL constructs a full URL to n image.
// It returns nil if BASE_URL is not set.
func GetImageURL(path string) *string {
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
