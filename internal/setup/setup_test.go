package setup

import (
	"log/slog"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	status := Check(logger)

	// We can't guarantee tshark is installed, but the function shouldn't panic.
	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if status.Installed {
		if status.TSharkPath == "" {
			t.Error("installed=true but TSharkPath is empty")
		}
		if status.Version == "" {
			t.Error("installed=true but Version is empty")
		}
		t.Logf("tshark found: %s (version %s)", status.TSharkPath, status.Version)
	} else {
		t.Log("tshark not installed on this system")
	}
}

func TestPrintInstructions(t *testing.T) {
	instructions := PrintInstructions()

	if instructions == "" {
		t.Fatal("expected non-empty instructions")
	}

	if !strings.Contains(instructions, "tshark.dev") {
		t.Error("expected reference to tshark.dev")
	}

	// Should contain platform-specific instructions.
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(instructions, "winget") && !strings.Contains(instructions, "choco") {
			t.Error("expected Windows-specific instructions")
		}
	case "darwin":
		if !strings.Contains(instructions, "brew") {
			t.Error("expected macOS-specific instructions")
		}
	case "linux":
		if !strings.Contains(instructions, "apt-get") {
			t.Error("expected Linux-specific instructions")
		}
	}
}

func TestFindBinary(t *testing.T) {
	// Test with a binary that definitely exists.
	var knownBinary string
	if runtime.GOOS == "windows" {
		knownBinary = "cmd"
	} else {
		knownBinary = "sh"
	}

	path, err := findBinary(knownBinary)
	if err != nil {
		t.Fatalf("expected to find %s: %v", knownBinary, err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}

	// Test with a binary that doesn't exist.
	_, err = findBinary("nonexistent_binary_xyz_123")
	if err == nil {
		t.Error("expected error for nonexistent binary")
	}
}
