package utils

import (
	"regexp"
	"strings"
)

var (
	// Regular expression to match multiple whitespace characters
	whitespaceRegex = regexp.MustCompile(`\s+`)
)

// NormalizeContent normalizes input content to ensure consistent cache hits
// and moderation API results for semantically identical content.
// Normalization includes:
// - Trimming whitespace
// - Replacing multiple whitespace characters with a single space
// - Converting to lowercase
func NormalizeContent(content string) string {
	// Trim leading and trailing whitespace
	content = strings.TrimSpace(content)
	
	// Replace multiple whitespace characters with a single space
	content = whitespaceRegex.ReplaceAllString(content, " ")
	
	// Convert to lowercase for case-insensitive matching
	content = strings.ToLower(content)
	
	return content
}
