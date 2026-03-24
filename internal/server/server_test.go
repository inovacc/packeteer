package server

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/inovacc/packeteer/internal/executor"
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
