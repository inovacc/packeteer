package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupResourceServer(t *testing.T, mock *executor.MockExecutor, captureDir string) *mcp.ClientSession {
	t.Helper()
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0"}, nil)
	RegisterResources(server, mock, captureDir)

	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, st, nil); err != nil {
		t.Fatalf("connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "tc", Version: "1.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func TestRegisterResources_ListCaptures(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "a.pcap"), []byte("fake"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.pcapng"), []byte("fake"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "skip.txt"), []byte("not pcap"), 0644)

	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, tmpDir)
	ctx := context.Background()

	result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "sharkline://captures"})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	text := result.Contents[0].Text
	if !strings.Contains(text, "a.pcap") {
		t.Error("expected a.pcap")
	}
	if !strings.Contains(text, "b.pcapng") {
		t.Error("expected b.pcapng")
	}
	if strings.Contains(text, "skip.txt") {
		t.Error("should not include .txt files")
	}
}

func TestRegisterResources_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, tmpDir)

	result, _ := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: "sharkline://captures"})
	if !strings.Contains(result.Contents[0].Text, "No capture files") {
		t.Error("expected empty message")
	}
}

func TestRegisterResources_NoCaptureDir(t *testing.T) {
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, "")

	result, _ := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{URI: "sharkline://captures"})
	if !strings.Contains(result.Contents[0].Text, "No captures directory") {
		t.Error("expected no-dir message")
	}
}

func TestRegisterResources_ReadSpecificCapture(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.pcap"), []byte("pcap data"), 0644)

	mock := executor.NewMockExecutor()
	mock.Responses["capinfos"] = executor.MockResponse{
		Stdout: []byte("File: test.pcap\nPackets: 42\n"),
	}
	cs := setupResourceServer(t, mock, tmpDir)

	result, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/test.pcap",
	})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(result.Contents[0].Text, "Packets: 42") {
		t.Error("expected capinfos output")
	}
}

func TestRegisterResources_InvalidExtension(t *testing.T) {
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, t.TempDir())

	_, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/evil.exe",
	})
	if err == nil {
		t.Error("expected error for .exe")
	}
}

func TestRegisterResources_PathTraversal(t *testing.T) {
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, t.TempDir())

	_, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/../../../etc/passwd.pcap",
	})
	if err == nil {
		t.Error("expected error for path traversal")
	}
}

func TestRegisterResources_EmptyFilename(t *testing.T) {
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, t.TempDir())

	_, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/",
	})
	if err == nil {
		t.Error("expected error for empty filename")
	}
}

func TestRegisterResources_NonexistentFile(t *testing.T) {
	mock := executor.NewMockExecutor()
	cs := setupResourceServer(t, mock, t.TempDir())

	_, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/nofile.pcap",
	})
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestRegisterResources_CapinfosError(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "broken.pcap"), []byte("data"), 0644)

	mock := executor.NewMockExecutor()
	mock.Responses["capinfos"] = executor.MockResponse{
		Err: os.ErrPermission,
	}
	cs := setupResourceServer(t, mock, tmpDir)

	result, err := cs.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: "sharkline://captures/broken.pcap",
	})
	if err != nil {
		t.Fatalf("should return result with error text, not error: %v", err)
	}
	if !strings.Contains(result.Contents[0].Text, "Error reading") {
		t.Error("expected error message in content")
	}
}
