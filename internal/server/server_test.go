package server

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNew_RegistersAllTools(t *testing.T) {
	mock := executor.NewMockExecutor()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	srv := New(mock, logger)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()

	_, err := srv.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer func() { _ = cs.Close() }()

	result, err := cs.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}

	expectedTools := map[string]bool{
		"list_interfaces":  false,
		"capture_packets":  false,
		"read_pcap":        false,
		"extract_fields":   false,
		"get_statistics":   false,
		"get_capture_info": false,
		"filter_pcap":      false,
		"merge_pcaps":      false,
		"list_protocols":   false,
		"decode_packet":    false,
	}

	for _, tool := range result.Tools {
		if _, ok := expectedTools[tool.Name]; ok {
			expectedTools[tool.Name] = true
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("tool %q not registered", name)
		}
	}

	if len(result.Tools) != 10 {
		t.Errorf("expected 10 tools, got %d", len(result.Tools))
	}
}

func TestNew_CallTool(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("1. eth0\n2. lo (Loopback)\n"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	srv := New(mock, logger)

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()

	_, err := srv.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer func() { _ = cs.Close() }()

	result, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "list_interfaces",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}

	if result.IsError {
		t.Fatal("expected success from list_interfaces")
	}

	if len(result.Content) == 0 {
		t.Fatal("expected non-empty content")
	}

	text := result.Content[0].(*mcp.TextContent).Text
	if text == "" {
		t.Fatal("expected non-empty text content")
	}
}

// helper to create a connected client session.
func setupClient(t *testing.T, mock *executor.MockExecutor, opts ...Option) *mcp.ClientSession {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	srv := New(mock, logger, opts...)

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()

	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}

	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func TestPrompts_ListAndGet(t *testing.T) {
	cs := setupClient(t, executor.NewMockExecutor())
	ctx := context.Background()

	t.Run("list prompts", func(t *testing.T) {
		result, err := cs.ListPrompts(ctx, &mcp.ListPromptsParams{})
		if err != nil {
			t.Fatalf("list prompts failed: %v", err)
		}

		expectedPrompts := map[string]bool{
			"analyze-traffic":        false,
			"investigate-connection": false,
			"security-scan":          false,
		}

		for _, p := range result.Prompts {
			if _, ok := expectedPrompts[p.Name]; ok {
				expectedPrompts[p.Name] = true
			}
		}

		for name, found := range expectedPrompts {
			if !found {
				t.Errorf("prompt %q not registered", name)
			}
		}
	})

	t.Run("get analyze-traffic", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "analyze-traffic",
			Arguments: map[string]string{
				"file_path": "/tmp/test.pcap",
				"focus":     "dns",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		if len(result.Messages) == 0 {
			t.Fatal("expected messages")
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "test.pcap") {
			t.Error("expected file_path in prompt text")
		}
		if !strings.Contains(text, "dns") {
			t.Error("expected focus area in prompt text")
		}
	})

	t.Run("get analyze-traffic default focus", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "analyze-traffic",
			Arguments: map[string]string{
				"file_path": "/tmp/test.pcap",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "overview") {
			t.Error("expected default focus 'overview'")
		}
	})

	t.Run("get investigate-connection", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "investigate-connection",
			Arguments: map[string]string{
				"file_path": "/tmp/test.pcap",
				"source_ip": "192.168.1.1",
				"dest_ip":   "10.0.0.1",
				"port":      "443",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "192.168.1.1") {
			t.Error("expected source_ip in prompt")
		}
		if !strings.Contains(text, "tcp.port == 443") {
			t.Error("expected port filter in prompt")
		}
	})

	t.Run("get investigate-connection no port", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "investigate-connection",
			Arguments: map[string]string{
				"file_path": "/tmp/test.pcap",
				"source_ip": "192.168.1.1",
				"dest_ip":   "10.0.0.1",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if strings.Contains(text, "tcp.port") {
			t.Error("should not include port filter when port is empty")
		}
	})

	t.Run("get security-scan", func(t *testing.T) {
		result, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "security-scan",
			Arguments: map[string]string{
				"file_path": "/tmp/test.pcap",
			},
		})
		if err != nil {
			t.Fatalf("get prompt failed: %v", err)
		}
		text := result.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "DNS exfiltration") {
			t.Error("expected security checks in prompt")
		}
	})
}

func TestResources_ListAndRead(t *testing.T) {
	// Create a temp dir with a fake pcap file for the resource listing.
	tmpDir := t.TempDir()
	fakePcap := filepath.Join(tmpDir, "sample.pcap")
	if err := os.WriteFile(fakePcap, []byte("fake pcap data for testing"), 0644); err != nil {
		t.Fatalf("failed to create fake pcap: %v", err)
	}

	mock := executor.NewMockExecutor()
	mock.Responses["capinfos"] = executor.MockResponse{
		Stdout: []byte("File name: sample.pcap\nPackets: 100\nDuration: 5.0 seconds\n"),
	}

	cs := setupClient(t, mock, WithCaptureDir(tmpDir))
	ctx := context.Background()

	t.Run("list resources", func(t *testing.T) {
		result, err := cs.ListResources(ctx, &mcp.ListResourcesParams{})
		if err != nil {
			t.Fatalf("list resources failed: %v", err)
		}
		found := false
		for _, r := range result.Resources {
			if r.URI == "sharkline://captures" {
				found = true
			}
		}
		if !found {
			t.Error("expected sharkline://captures resource")
		}
	})

	t.Run("list resource templates", func(t *testing.T) {
		result, err := cs.ListResourceTemplates(ctx, &mcp.ListResourceTemplatesParams{})
		if err != nil {
			t.Fatalf("list resource templates failed: %v", err)
		}
		found := false
		for _, rt := range result.ResourceTemplates {
			if strings.Contains(rt.URITemplate, "captures") {
				found = true
			}
		}
		if !found {
			t.Error("expected captures resource template")
		}
	})

	t.Run("read captures list", func(t *testing.T) {
		result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "sharkline://captures",
		})
		if err != nil {
			t.Fatalf("read resource failed: %v", err)
		}
		if len(result.Contents) == 0 {
			t.Fatal("expected contents")
		}
		text := result.Contents[0].Text
		if !strings.Contains(text, "sample.pcap") {
			t.Errorf("expected sample.pcap in listing, got: %s", text)
		}
	})

	t.Run("read specific capture", func(t *testing.T) {
		result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "sharkline://captures/sample.pcap",
		})
		if err != nil {
			t.Fatalf("read resource failed: %v", err)
		}
		if len(result.Contents) == 0 {
			t.Fatal("expected contents")
		}
		text := result.Contents[0].Text
		if !strings.Contains(text, "Packets: 100") {
			t.Errorf("expected capinfos output, got: %s", text)
		}
	})

	t.Run("read invalid extension", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "sharkline://captures/malware.exe",
		})
		if err == nil {
			t.Fatal("expected error for invalid extension")
		}
	})

	t.Run("read nonexistent file", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "sharkline://captures/nonexistent.pcap",
		})
		if err == nil {
			t.Fatal("expected error for nonexistent file")
		}
	})
}

func TestResources_NoCaptureDir(t *testing.T) {
	cs := setupClient(t, executor.NewMockExecutor())
	ctx := context.Background()

	result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: "sharkline://captures",
	})
	if err != nil {
		t.Fatalf("read resource failed: %v", err)
	}
	text := result.Contents[0].Text
	if !strings.Contains(text, "No captures directory") {
		t.Errorf("expected no-dir message, got: %s", text)
	}
}
