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

func TestGetTSharkVersion(t *testing.T) {
	status := Check(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))
	if !status.Installed {
		t.Skip("tshark not installed")
	}

	version, err := getTSharkVersion(status.TSharkPath)
	if err != nil {
		t.Fatalf("getTSharkVersion failed: %v", err)
	}
	if version == "" {
		t.Error("expected non-empty version")
	}
	// Version should match X.Y.Z pattern.
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		t.Errorf("expected X.Y.Z format, got %q", version)
	}
	t.Logf("version: %s", version)
}

func TestPrintInstructions_AllPlatforms(t *testing.T) {
	// Just verify the function doesn't panic and returns content.
	instructions := PrintInstructions()
	if !strings.Contains(instructions, "Installation Guide") {
		t.Error("expected 'Installation Guide' header")
	}
	if !strings.Contains(instructions, "tshark.dev") {
		t.Error("expected tshark.dev reference")
	}
}

func TestCheck_StatusFields(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	status := Check(logger)

	if status.Installed {
		// All tools should be found if tshark is installed.
		if status.CapinfosPath == "" {
			t.Error("expected capinfos when tshark is installed")
		}
		if status.EditcapPath == "" {
			t.Error("expected editcap when tshark is installed")
		}
		if status.MergecapPath == "" {
			t.Error("expected mergecap when tshark is installed")
		}
	}
}

func TestFindBinary_PlatformPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Test that findBinary checks Program Files.
		path, err := findBinary("tshark")
		if err == nil {
			if !strings.Contains(path, "Wireshark") {
				t.Errorf("expected Wireshark in path, got %s", path)
			}
		}
	}
}
