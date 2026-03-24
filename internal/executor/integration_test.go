//go:build integration

package executor

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestRealExecutor_Integration_ListInterfaces(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	exec := NewRealExecutor(logger, "")

	// Verify tshark is available.
	if _, err := exec.BinaryPath("tshark"); err != nil {
		t.Skipf("tshark not available: %v", err)
	}

	stdout, stderr, err := exec.Execute(context.Background(), "tshark", []string{"-D"})
	if err != nil {
		// Npcap may not be installed — tshark -D fails without it but tshark still works.
		if strings.Contains(string(stderr), "Npcap") || strings.Contains(err.Error(), "Npcap") {
			t.Skipf("tshark works but Npcap not installed (needed for interface listing): %v", err)
		}
		t.Fatalf("tshark -D failed: %v", err)
	}

	output := string(stdout)
	if output == "" {
		t.Fatal("expected non-empty interface list")
	}

	// Should contain at least one numbered interface.
	if !strings.Contains(output, "1.") {
		t.Errorf("expected numbered interface list, got: %s", output)
	}

	t.Logf("interfaces:\n%s", output)
}

func TestRealExecutor_Integration_ListProtocols(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	exec := NewRealExecutor(logger, "")

	if _, err := exec.BinaryPath("tshark"); err != nil {
		t.Skipf("tshark not available: %v", err)
	}

	stdout, _, err := exec.Execute(context.Background(), "tshark", []string{"-G", "protocols"})
	if err != nil {
		t.Fatalf("tshark -G protocols failed: %v", err)
	}

	output := string(stdout)
	if !strings.Contains(strings.ToLower(output), "tcp") {
		t.Error("expected TCP in protocol list")
	}
	if !strings.Contains(strings.ToLower(output), "udp") {
		t.Error("expected UDP in protocol list")
	}
}

func TestRealExecutor_Integration_TsharkVersion(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	exec := NewRealExecutor(logger, "")

	if _, err := exec.BinaryPath("tshark"); err != nil {
		t.Skipf("tshark not available: %v", err)
	}

	stdout, _, err := exec.Execute(context.Background(), "tshark", []string{"--version"})
	if err != nil {
		// tshark --version may exit non-zero on some versions, check stdout.
		if len(stdout) == 0 {
			t.Fatalf("tshark --version failed with no output: %v", err)
		}
	}

	output := string(stdout)
	if !strings.Contains(output, "TShark") && !strings.Contains(output, "Wireshark") {
		t.Errorf("expected TShark/Wireshark in version output, got: %s", output[:min(len(output), 200)])
	}

	t.Logf("version: %s", strings.Split(output, "\n")[0])
}

func TestRealExecutor_Integration_CapinfosAvailable(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	exec := NewRealExecutor(logger, "")

	for _, binary := range []string{"capinfos", "editcap", "mergecap"} {
		t.Run(binary, func(t *testing.T) {
			path, err := exec.BinaryPath(binary)
			if err != nil {
				t.Skipf("%s not available: %v", binary, err)
			}
			t.Logf("%s found at: %s", binary, path)
		})
	}
}
