package output

import (
	"encoding/json"
	"testing"
)

func TestParseFieldOutput(t *testing.T) {
	t.Run("valid tab-separated", func(t *testing.T) {
		input := "192.168.1.1\t93.184.216.34\t80\n10.0.0.1\t8.8.8.8\t443\n"
		fields := []string{"ip.src", "ip.dst", "tcp.dstport"}

		result, err := ParseFieldOutput([]byte(input), fields)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed FieldsResult
		json.Unmarshal(result, &parsed)

		if parsed.Total != 2 {
			t.Errorf("expected 2 rows, got %d", parsed.Total)
		}
		if parsed.Rows[0].Fields["ip.src"] != "192.168.1.1" {
			t.Errorf("expected 192.168.1.1, got %s", parsed.Rows[0].Fields["ip.src"])
		}
		if parsed.Rows[1].Fields["tcp.dstport"] != "443" {
			t.Errorf("expected 443, got %s", parsed.Rows[1].Fields["tcp.dstport"])
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result, _ := ParseFieldOutput([]byte(""), []string{"ip.src"})
		var parsed FieldsResult
		json.Unmarshal(result, &parsed)
		if parsed.Total != 0 {
			t.Errorf("expected 0 rows, got %d", parsed.Total)
		}
	})

	t.Run("fewer values than fields", func(t *testing.T) {
		input := "192.168.1.1\n"
		fields := []string{"ip.src", "ip.dst", "tcp.port"}

		result, _ := ParseFieldOutput([]byte(input), fields)
		var parsed FieldsResult
		json.Unmarshal(result, &parsed)

		if parsed.Rows[0].Fields["ip.dst"] != "" {
			t.Error("expected empty for missing field")
		}
	})

	t.Run("preserves field names", func(t *testing.T) {
		result, _ := ParseFieldOutput([]byte("a\tb\n"), []string{"x", "y"})
		var parsed FieldsResult
		json.Unmarshal(result, &parsed)

		if len(parsed.FieldNames) != 2 || parsed.FieldNames[0] != "x" {
			t.Error("expected field names preserved")
		}
	})
}
