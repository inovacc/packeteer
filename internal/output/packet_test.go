package output

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseTSharkJSON(t *testing.T) {
	t.Run("valid json array", func(t *testing.T) {
		input := `[
			{
				"_index": "packets-1",
				"_source": {
					"layers": {
						"frame": {
							"frame.time_relative": "0.000000",
							"frame.len": "66",
							"frame.protocols": "eth:ethertype:ip:tcp"
						},
						"ip": {
							"ip.src": "192.168.1.100",
							"ip.dst": "93.184.216.34"
						},
						"tcp": {
							"tcp.srcport": "49152",
							"tcp.dstport": "80",
							"tcp.flags.str": "SYN"
						}
					}
				}
			},
			{
				"_index": "packets-2",
				"_source": {
					"layers": {
						"frame": {
							"frame.time_relative": "0.050000",
							"frame.len": "60",
							"frame.protocols": "eth:ethertype:ip:udp:dns"
						},
						"ip": {
							"ip.src": "192.168.1.100",
							"ip.dst": "8.8.8.8"
						},
						"dns": {
							"dns.qry.name": "example.com"
						}
					}
				}
			}
		]`

		result, err := ParseTSharkJSON([]byte(input), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed ParseResult
		if err := json.Unmarshal(result, &parsed); err != nil {
			t.Fatalf("failed to parse result: %v", err)
		}

		if parsed.Total != 2 {
			t.Errorf("expected 2 packets, got %d", parsed.Total)
		}

		if parsed.Packets[0].Source != "192.168.1.100" {
			t.Errorf("expected source 192.168.1.100, got %s", parsed.Packets[0].Source)
		}
		if parsed.Packets[0].Dest != "93.184.216.34" {
			t.Errorf("expected dest 93.184.216.34, got %s", parsed.Packets[0].Dest)
		}
		if !strings.Contains(parsed.Packets[0].Info, "TCP") {
			t.Errorf("expected TCP info, got %s", parsed.Packets[0].Info)
		}
		if !strings.Contains(parsed.Packets[1].Info, "DNS") {
			t.Errorf("expected DNS info, got %s", parsed.Packets[1].Info)
		}
	})

	t.Run("truncation", func(t *testing.T) {
		input := `[{"_source":{"layers":{}}},{"_source":{"layers":{}}},{"_source":{"layers":{}}}]`
		result, err := ParseTSharkJSON([]byte(input), 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var parsed ParseResult
		json.Unmarshal(result, &parsed)

		if len(parsed.Packets) != 2 {
			t.Errorf("expected 2 packets after truncation, got %d", len(parsed.Packets))
		}
		if parsed.Total != 3 {
			t.Errorf("expected total 3, got %d", parsed.Total)
		}
		if !parsed.Truncated {
			t.Error("expected truncated=true")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := ParseTSharkJSON([]byte(""), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var parsed ParseResult
		json.Unmarshal(result, &parsed)
		if parsed.Total != 0 {
			t.Errorf("expected 0 packets, got %d", parsed.Total)
		}
	})

	t.Run("invalid json falls back", func(t *testing.T) {
		input := "not json at all"
		result, err := ParseTSharkJSON([]byte(input), 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(result) != input {
			t.Error("expected original data returned on parse failure")
		}
	})

	t.Run("http packet info", func(t *testing.T) {
		input := `[{"_source":{"layers":{"frame":{},"ip":{"ip.src":"1.2.3.4","ip.dst":"5.6.7.8"},"http":{"http.request.method":"GET","http.request.uri":"/api/data"}}}}]`
		result, _ := ParseTSharkJSON([]byte(input), 0)
		var parsed ParseResult
		json.Unmarshal(result, &parsed)
		if !strings.Contains(parsed.Packets[0].Info, "HTTP GET /api/data") {
			t.Errorf("expected HTTP info, got %s", parsed.Packets[0].Info)
		}
	})

	t.Run("tls packet info", func(t *testing.T) {
		input := `[{"_source":{"layers":{"frame":{},"ip":{"ip.src":"1.2.3.4","ip.dst":"5.6.7.8"},"tls":{"tls.handshake.extensions_server_name":"example.com"}}}}]`
		result, _ := ParseTSharkJSON([]byte(input), 0)
		var parsed ParseResult
		json.Unmarshal(result, &parsed)
		if !strings.Contains(parsed.Packets[0].Info, "TLS") {
			t.Errorf("expected TLS info, got %s", parsed.Packets[0].Info)
		}
	})
}
