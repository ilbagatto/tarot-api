package utils

import "strings"

// ParseBoolParam converts query string values into boolean
func ParseBoolParam(param string) bool {
	switch strings.ToLower(param) {
	case "true", "1", "yes":
		return true
	default:
		return false
	}
}
