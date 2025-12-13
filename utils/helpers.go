package utils

import (
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
