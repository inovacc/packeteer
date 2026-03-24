package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Packet represents a parsed tshark JSON packet.
type Packet struct {
	Index  string                 `json:"_index,omitempty"`
	Source PacketSource           `json:"_source,omitempty"`
	Raw    map[string]interface{} `json:"-"`
}

// PacketSource contains the layers of a packet.
type PacketSource struct {
	Layers map[string]interface{} `json:"layers,omitempty"`
}

// PacketSummary is a simplified view of a packet for AI consumption.
type PacketSummary struct {
	Number   int    `json:"number,omitempty"`
	Time     string `json:"time,omitempty"`
	Source   string `json:"source,omitempty"`
	Dest     string `json:"destination,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Length   int    `json:"length,omitempty"`
	Info     string `json:"info,omitempty"`
}

// ParseResult holds parsed tshark JSON output with metadata.
type ParseResult struct {
	Packets   []PacketSummary `json:"packets"`
	Total     int             `json:"total"`
	Truncated bool            `json:"truncated"`
}

// ParseTSharkJSON parses tshark -T json output into structured packet summaries.
// Returns the original output if parsing fails (graceful fallback).
func ParseTSharkJSON(data []byte, maxPackets int) ([]byte, error) {
	data = []byte(strings.TrimSpace(string(data)))
	if len(data) == 0 {
		return json.Marshal(ParseResult{Packets: []PacketSummary{}, Total: 0})
	}

	var rawPackets []map[string]interface{}
	if err := json.Unmarshal(data, &rawPackets); err != nil {
		// Not valid JSON array — return original data as-is.
		return data, nil
	}

	total := len(rawPackets)
	if maxPackets > 0 && total > maxPackets {
		rawPackets = rawPackets[:maxPackets]
	}

	summaries := make([]PacketSummary, 0, len(rawPackets))
	for i, raw := range rawPackets {
		summary := extractSummary(raw, i+1)
		summaries = append(summaries, summary)
	}

	result := ParseResult{
		Packets:   summaries,
		Total:     total,
		Truncated: maxPackets > 0 && total > maxPackets,
	}

	return json.MarshalIndent(result, "", "  ")
}

func extractSummary(raw map[string]interface{}, index int) PacketSummary {
	s := PacketSummary{Number: index}

	source, ok := raw["_source"].(map[string]interface{})
	if !ok {
		return s
	}

	layers, ok := source["layers"].(map[string]interface{})
	if !ok {
		return s
	}

	// Extract frame info.
	if frame, ok := layers["frame"].(map[string]interface{}); ok {
		s.Time = getStr(frame, "frame.time_relative")
		if l, ok := frame["frame.len"]; ok {
			s.Length = getInt(l)
		}
		s.Protocol = getStr(frame, "frame.protocols")
	}

	// Extract IP info.
	if ip, ok := layers["ip"].(map[string]interface{}); ok {
		s.Source = getStr(ip, "ip.src")
		s.Dest = getStr(ip, "ip.dst")
	} else if ipv6, ok := layers["ipv6"].(map[string]interface{}); ok {
		s.Source = getStr(ipv6, "ipv6.src")
		s.Dest = getStr(ipv6, "ipv6.dst")
	}

	// Build info string from highest-layer protocol.
	s.Info = buildInfo(layers)

	return s
}

func buildInfo(layers map[string]interface{}) string {
	// Check for common application protocols in order.
	for _, proto := range []string{"http", "dns", "tls", "tcp", "udp"} {
		if layer, ok := layers[proto].(map[string]interface{}); ok {
			switch proto {
			case "http":
				method := getStr(layer, "http.request.method")
				uri := getStr(layer, "http.request.uri")
				code := getStr(layer, "http.response.code")
				if method != "" {
					return fmt.Sprintf("HTTP %s %s", method, uri)
				}
				if code != "" {
					return fmt.Sprintf("HTTP %s", code)
				}
			case "dns":
				qname := getStr(layer, "dns.qry.name")
				if qname != "" {
					return fmt.Sprintf("DNS %s", qname)
				}
			case "tls":
				sni := getStr(layer, "tls.handshake.extensions_server_name")
				if sni != "" {
					return fmt.Sprintf("TLS → %s", sni)
				}
			case "tcp":
				flags := getStr(layer, "tcp.flags.str")
				srcPort := getStr(layer, "tcp.srcport")
				dstPort := getStr(layer, "tcp.dstport")
				if flags != "" {
					return fmt.Sprintf("TCP %s→%s [%s]", srcPort, dstPort, flags)
				}
				if srcPort != "" {
					return fmt.Sprintf("TCP %s→%s", srcPort, dstPort)
				}
			case "udp":
				srcPort := getStr(layer, "udp.srcport")
				dstPort := getStr(layer, "udp.dstport")
				if srcPort != "" {
					return fmt.Sprintf("UDP %s→%s", srcPort, dstPort)
				}
			}
		}
	}
	return ""
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		var n int
		fmt.Sscanf(val, "%d", &n)
		return n
	default:
		return 0
	}
}
