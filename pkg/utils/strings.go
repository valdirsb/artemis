package utils

import (
	"strings"
)

// TrimSpaces removes leading and trailing spaces from a string.
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

// ToLower converts a string to lowercase.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts a string to uppercase.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Split separates a string into substrings based on the specified delimiter.
func Split(s, delimiter string) []string {
	return strings.Split(s, delimiter)
}

// Join concatenates the elements of a slice into a single string with the specified delimiter.
func Join(elements []string, delimiter string) string {
	return strings.Join(elements, delimiter)
}

// Contains checks if a substring is present within a string.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Replace replaces occurrences of old with new in a string.
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}