package output

import (
	"fmt"
)

const (
	DefaultMaxBytes = 512 * 1024 // 512KB default
	TruncationNote  = "\n\n[Output truncated: %d bytes shown of %d total. Use filters to narrow results.]"
)

// Truncate limits output to maxBytes, appending a truncation notice if needed.
func Truncate(data []byte, maxBytes int) ([]byte, bool) {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}

	if len(data) <= maxBytes {
		return data, false
	}

	notice := fmt.Sprintf(TruncationNote, maxBytes, len(data))
	result := make([]byte, 0, maxBytes+len(notice))
	result = append(result, data[:maxBytes]...)
	result = append(result, []byte(notice)...)
	return result, true
}

// FormatResult creates a structured text result with metadata.
func FormatResult(content string, metadata map[string]string) string {
	if len(metadata) == 0 {
		return content
	}

	header := ""
	for k, v := range metadata {
		header += fmt.Sprintf("%s: %s\n", k, v)
	}
	header += "---\n"

	return header + content
}
