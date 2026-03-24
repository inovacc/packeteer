package tools

import (
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name string
		b    int64
		want string
	}{
		{"bytes", 512, "512 B"},
		{"zero", 0, "0 B"},
		{"kilobytes", 2048, "2.0 KB"},
		{"megabytes", 5 * 1024 * 1024, "5.0 MB"},
		{"gigabytes", 3 * 1024 * 1024 * 1024, "3.0 GB"},
		{"fractional KB", 1536, "1.5 KB"},
		{"fractional MB", 1572864, "1.5 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSize(tt.b)
			if got != tt.want {
				t.Errorf("formatSize(%d) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
