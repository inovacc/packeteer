package output

import (
	"strings"
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		maxBytes     int
		wantTrunc    bool
		wantContains string
	}{
		{
			name:      "under limit",
			data:      []byte("short"),
			maxBytes:  100,
			wantTrunc: false,
		},
		{
			name:      "exact limit",
			data:      []byte("12345"),
			maxBytes:  5,
			wantTrunc: false,
		},
		{
			name:         "over limit",
			data:         []byte("this is a longer string that exceeds the limit"),
			maxBytes:     10,
			wantTrunc:    true,
			wantContains: "Output truncated",
		},
		{
			name:      "zero maxBytes uses default",
			data:      []byte("short"),
			maxBytes:  0,
			wantTrunc: false,
		},
		{
			name:      "negative maxBytes uses default",
			data:      []byte("short"),
			maxBytes:  -1,
			wantTrunc: false,
		},
		{
			name:      "empty data",
			data:      []byte{},
			maxBytes:  100,
			wantTrunc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, truncated := Truncate(tt.data, tt.maxBytes)
			if truncated != tt.wantTrunc {
				t.Errorf("truncated = %v, want %v", truncated, tt.wantTrunc)
			}
			if tt.wantContains != "" && !strings.Contains(string(result), tt.wantContains) {
				t.Errorf("result should contain %q, got %q", tt.wantContains, result)
			}
			if !tt.wantTrunc && string(result) != string(tt.data) {
				t.Errorf("non-truncated result should equal input")
			}
		})
	}
}

func TestTruncate_PreservesPrefix(t *testing.T) {
	data := []byte("ABCDEFGHIJ")
	result, truncated := Truncate(data, 5)
	if !truncated {
		t.Fatal("expected truncation")
	}
	if !strings.HasPrefix(string(result), "ABCDE") {
		t.Errorf("should preserve first 5 bytes, got %q", result)
	}
}

func TestFormatResult(t *testing.T) {
	t.Run("with metadata", func(t *testing.T) {
		result := FormatResult("content here", map[string]string{
			"File": "test.pcap",
		})
		if !strings.Contains(result, "File: test.pcap") {
			t.Error("should contain metadata")
		}
		if !strings.Contains(result, "---") {
			t.Error("should contain separator")
		}
		if !strings.Contains(result, "content here") {
			t.Error("should contain content")
		}
	})

	t.Run("empty metadata", func(t *testing.T) {
		result := FormatResult("content", map[string]string{})
		if result != "content" {
			t.Errorf("empty metadata should return content only, got %q", result)
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		result := FormatResult("content", nil)
		if result != "content" {
			t.Errorf("nil metadata should return content only, got %q", result)
		}
	})
}
