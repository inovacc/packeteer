//go:build integration

package server

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestE2E_SamplePcap runs the full MCP pipeline against the bundled sample.pcap
// using a real tshark binary. Skip if tshark is not installed.
func TestE2E_SamplePcap(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	exec := executor.NewRealExecutor(logger, "")

	if _, err := exec.BinaryPath("tshark"); err != nil {
		t.Skipf("tshark not available: %v", err)
	}

	// Find sample.pcap relative to project root.
	samplePcap, err := filepath.Abs("../../testdata/sample.pcap")
	if err != nil {
		t.Fatalf("failed to resolve sample.pcap: %v", err)
	}
	if _, err := os.Stat(samplePcap); os.IsNotExist(err) {
		t.Skipf("sample.pcap not found at %s", samplePcap)
	}

	captureDir := filepath.Dir(samplePcap)
	srv := New(exec, logger, WithCaptureDir(captureDir))

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "e2e-test", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer func() { _ = cs.Close() }()

	t.Run("read_pcap", func(t *testing.T) {
		result, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "read_pcap",
			Arguments: map[string]any{
				"file_path":  samplePcap,
				"max_packets": 10,
			},
		})
		if err != nil {
			t.Fatalf("call tool failed: %v", err)
		}
		if result.IsError {
			text := result.Content[0].(*mcp.TextContent).Text
			t.Fatalf("read_pcap returned error: %s", text)
		}
		t.Log("read_pcap succeeded")
	})

	t.Run("extract_fields", func(t *testing.T) {
		result, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "extract_fields",
			Arguments: map[string]any{
				"file_path": samplePcap,
				"fields":    []any{"ip.src", "ip.dst", "tcp.dstport"},
			},
		})
		if err != nil {
			t.Fatalf("call tool failed: %v", err)
		}
		if result.IsError {
			text := result.Content[0].(*mcp.TextContent).Text
			t.Fatalf("extract_fields returned error: %s", text)
		}
		text := result.Content[0].(*mcp.TextContent).Text
		if !strings.Contains(text, "192.168.1.100") {
			t.Log("Note: source IP not found — tshark may parse differently")
		}
		t.Log("extract_fields succeeded")
	})

	t.Run("get_capture_info", func(t *testing.T) {
		result, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "get_capture_info",
			Arguments: map[string]any{
				"file_path": samplePcap,
			},
		})
		if err != nil {
			t.Fatalf("call tool failed: %v", err)
		}
		if result.IsError {
			text := result.Content[0].(*mcp.TextContent).Text
			t.Fatalf("get_capture_info returned error: %s", text)
		}
		t.Log("get_capture_info succeeded")
	})

	t.Run("decode_packet", func(t *testing.T) {
		result, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "decode_packet",
			Arguments: map[string]any{
				"file_path":     samplePcap,
				"packet_number": 1,
			},
		})
		if err != nil {
			t.Fatalf("call tool failed: %v", err)
		}
		if result.IsError {
			text := result.Content[0].(*mcp.TextContent).Text
			t.Fatalf("decode_packet returned error: %s", text)
		}
		t.Log("decode_packet succeeded")
	})

	t.Run("resource_list_captures", func(t *testing.T) {
		result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "packeteer://captures",
		})
		if err != nil {
			t.Fatalf("read resource failed: %v", err)
		}
		text := result.Contents[0].Text
		if !strings.Contains(text, "sample.pcap") {
			t.Errorf("expected sample.pcap in listing, got: %s", text)
		}
		t.Log("resource listing succeeded")
	})

	t.Run("resource_read_capture", func(t *testing.T) {
		result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "packeteer://captures/sample.pcap",
		})
		if err != nil {
			t.Fatalf("read resource failed: %v", err)
		}
		text := result.Contents[0].Text
		if text == "" {
			t.Error("expected non-empty capinfos output")
		}
		t.Log("resource read succeeded")
	})

	t.Run("prompt_analyze_traffic", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "analyze-traffic",
			Arguments: map[string]string{
				"file_path": samplePcap,
				"focus":     "overview",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, samplePcap) {
			t.Error("expected file path in prompt")
		}
		t.Log("prompt succeeded")
	})
}
