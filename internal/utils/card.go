package utils

import (
	"os"
	"strings"
)

// GetImageURL constructs a full URL to n image.
// It returns nil if STATIC_URL is not set.
func GetImageURL(path string, small bool) *string {
	base := os.Getenv("STATIC_URL")
	if base == "" {
		return nil
	}

	var imgPath string
	if small {
		imgPath = base + "/thumbnails"
	} else {
		imgPath = base + "/images"
	}

	var url string
	if strings.HasPrefix(path, "/") {
		url = imgPath + path
	} else {
		url = imgPath + "/" + path
	}

	return &url
}
