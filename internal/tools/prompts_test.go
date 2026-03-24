package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterPrompts(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0"}, nil)
	RegisterPrompts(server)

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, st, nil); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "tc", Version: "1.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}
	defer func() { _ = cs.Close() }()

	t.Run("lists 3 prompts", func(t *testing.T) {
		result, err := cs.ListPrompts(ctx, &mcp.ListPromptsParams{})
		if err != nil {
			t.Fatalf("list prompts: %v", err)
		}
		if len(result.Prompts) != 3 {
			t.Errorf("expected 3 prompts, got %d", len(result.Prompts))
		}
	})

	t.Run("analyze-traffic with focus", func(t *testing.T) {
		r, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "analyze-traffic",
			Arguments: map[string]string{"file_path": "/tmp/t.pcap", "focus": "http"},
		})
		if err != nil {
			t.Fatalf("get prompt: %v", err)
		}
		text := r.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "http") {
			t.Error("expected http focus")
		}
	})

	t.Run("analyze-traffic default focus", func(t *testing.T) {
		r, _ := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "analyze-traffic",
			Arguments: map[string]string{"file_path": "/tmp/t.pcap"},
		})
		text := r.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "overview") {
			t.Error("expected default overview focus")
		}
	})

	t.Run("investigate-connection with port", func(t *testing.T) {
		r, _ := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "investigate-connection",
			Arguments: map[string]string{
				"file_path": "/tmp/t.pcap", "source_ip": "1.2.3.4",
				"dest_ip": "5.6.7.8", "port": "443",
			},
		})
		text := r.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "tcp.port == 443") {
			t.Error("expected port filter")
		}
	})

	t.Run("investigate-connection without port", func(t *testing.T) {
		r, _ := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name: "investigate-connection",
			Arguments: map[string]string{
				"file_path": "/tmp/t.pcap", "source_ip": "1.2.3.4", "dest_ip": "5.6.7.8",
			},
		})
		text := r.Messages[0].Content.(*mcp.TextContent).Text
		if strings.Contains(text, "tcp.port") {
			t.Error("should not include port filter")
		}
	})

	t.Run("security-scan", func(t *testing.T) {
		r, _ := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "security-scan",
			Arguments: map[string]string{"file_path": "/tmp/t.pcap"},
		})
		text := r.Messages[0].Content.(*mcp.TextContent).Text
		if !strings.Contains(text, "DNS exfiltration") {
			t.Error("expected security checks")
		}
	})
}
