package utils

import (
	"os"
	"strings"
)

// TruncateString truncates a string to the specified length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// IsValidStatus checks if the manga status is valid
func IsValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"ongoing":   true,
		"completed": true,
		"hiatus":    true,
		"cancelled": true,
	}
	return validStatuses[status]
}

// MustReadFileBytes reads a file into a byte slice or panics
func MustReadFileBytes(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}
