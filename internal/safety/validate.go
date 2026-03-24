package safety

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	MaxCaptureTimeout = 30 * time.Second
	MaxCommandTimeout = 60 * time.Second
	MaxPacketCount    = 1000
	MaxOutputBytes    = 512 * 1024 // 512KB
)

var (
	allowedExtensions = map[string]bool{
		".pcap":   true,
		".pcapng": true,
		".cap":    true,
		".gz":     true,
	}

	// shellMetachars matches characters that could enable command injection.
	shellMetachars = regexp.MustCompile("[;|&$`\\\\\"'<>(){}\\[\\]!#~]")

	// validFilterChars allows alphanumeric, spaces, dots, colons, slashes,
	// comparison operators, parentheses for grouping, and common filter syntax.
	validFilterPattern = regexp.MustCompile(`^[a-zA-Z0-9\s._:,/=<>!()&|*?\-\[\]"']+$`)
)

// ValidateFilePath checks that a file path is safe for use with Wireshark tools.
func ValidateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("file path is required")
	}

	// Check for traversal in the raw path before cleaning resolves it away.
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	cleaned := filepath.Clean(path)

	ext := strings.ToLower(filepath.Ext(cleaned))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported file extension %q: allowed extensions are .pcap, .pcapng, .cap, .gz", ext)
	}

	return nil
}

// ValidateOutputPath checks that an output path is safe for writing.
func ValidateOutputPath(path string) error {
	if path == "" {
		return nil // optional
	}

	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed: %s", path)
	}

	cleaned := filepath.Clean(path)

	ext := strings.ToLower(filepath.Ext(cleaned))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported output extension %q: allowed extensions are .pcap, .pcapng, .cap, .gz", ext)
	}

	return nil
}

// SanitizeDisplayFilter validates a Wireshark display filter expression.
func SanitizeDisplayFilter(filter string) error {
	if filter == "" {
		return nil
	}

	if shellMetachars.MatchString(filter) {
		// Display filters can legitimately contain some of these chars,
		// so we check more specifically.
		if strings.ContainsAny(filter, ";`$\\") {
			return fmt.Errorf("display filter contains dangerous characters: %s", filter)
		}
	}

	return nil
}

// SanitizeCaptureFilter validates a BPF capture filter expression.
func SanitizeCaptureFilter(filter string) error {
	if filter == "" {
		return nil
	}

	if strings.ContainsAny(filter, ";`$\\{}") {
		return fmt.Errorf("capture filter contains dangerous characters: %s", filter)
	}

	return nil
}

// ClampTimeout returns the requested timeout clamped to the maximum.
func ClampTimeout(requested time.Duration, max time.Duration) time.Duration {
	if requested <= 0 {
		return max
	}
	if requested > max {
		return max
	}
	return requested
}

// ClampPacketCount returns the requested count clamped to the maximum.
func ClampPacketCount(requested int) int {
	if requested <= 0 {
		return MaxPacketCount
	}
	if requested > MaxPacketCount {
		return MaxPacketCount
	}
	return requested
}

// SanitizeInterfaceName validates a network interface name.
func SanitizeInterfaceName(name string) error {
	if name == "" {
		return fmt.Errorf("interface name is required")
	}

	if shellMetachars.MatchString(name) {
		return fmt.Errorf("interface name contains invalid characters: %s", name)
	}

	return nil
}

// SanitizeFieldName validates a tshark field name (e.g., "ip.src", "tcp.port").
func SanitizeFieldName(field string) error {
	if field == "" {
		return fmt.Errorf("field name is required")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, field)
	if !matched {
		return fmt.Errorf("invalid field name %q: only alphanumeric, dots, hyphens, and underscores allowed", field)
	}

	return nil
}

// SanitizeStatType validates a tshark statistics type (e.g., "io,phs", "conv,tcp").
func SanitizeStatType(stat string) error {
	if stat == "" {
		return fmt.Errorf("statistics type is required")
	}

	if !validFilterPattern.MatchString(stat) {
		return fmt.Errorf("invalid statistics type %q", stat)
	}

	return nil
}
