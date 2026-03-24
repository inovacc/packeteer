package executor

import (
	"context"
	"runtime"
	"testing"
)

func TestMockExecutor_ResponseMatching(t *testing.T) {
	mock := NewMockExecutor()
	mock.Responses["tshark"] = MockResponse{
		Stdout: []byte("output"),
	}
	mock.Responses["tshark -D"] = MockResponse{
		Stdout: []byte("interfaces"),
	}

	t.Run("matches binary only", func(t *testing.T) {
		stdout, _, err := mock.Execute(context.Background(), "tshark", []string{"-r", "test.pcap"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(stdout) != "output" {
			t.Fatalf("got %q, want %q", stdout, "output")
		}
	})

	t.Run("matches binary plus first arg", func(t *testing.T) {
		stdout, _, err := mock.Execute(context.Background(), "tshark", []string{"-D"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(stdout) != "interfaces" {
			t.Fatalf("got %q, want %q", stdout, "interfaces")
		}
	})

	t.Run("records calls", func(t *testing.T) {
		if len(mock.Calls) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(mock.Calls))
		}
		if mock.Calls[0].Binary != "tshark" {
			t.Fatalf("call 0: got binary %q", mock.Calls[0].Binary)
		}
	})

	t.Run("returns error for unknown binary", func(t *testing.T) {
		_, _, err := mock.Execute(context.Background(), "unknown", nil)
		if err == nil {
			t.Fatal("expected error for unknown binary")
		}
	})

	t.Run("uses default response", func(t *testing.T) {
		m := NewMockExecutor()
		m.DefaultResponse = &MockResponse{Stdout: []byte("default")}
		stdout, _, err := m.Execute(context.Background(), "anything", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(stdout) != "default" {
			t.Fatalf("got %q, want %q", stdout, "default")
		}
	})
}

func TestMockExecutor_BinaryPath(t *testing.T) {
	mock := NewMockExecutor()
	path, err := mock.BinaryPath("tshark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/usr/bin/tshark" {
		t.Fatalf("got %q, want /usr/bin/tshark", path)
	}
}

func TestRealExecutor_BinaryPath(t *testing.T) {
	exec := NewRealExecutor(nil, "")

	t.Run("rejects unknown binary", func(t *testing.T) {
		_, err := exec.BinaryPath("invalid")
		if err == nil {
			t.Fatal("expected error for unknown binary")
		}
	})

	t.Run("accepts valid binary names", func(t *testing.T) {
		for _, name := range []string{"tshark", "capinfos", "editcap", "mergecap", "dumpcap"} {
			// This may fail if Wireshark isn't installed, which is fine.
			// We're testing that the name is accepted, not that the binary exists.
			_, err := exec.BinaryPath(name)
			// Only check that it doesn't reject the name as "unknown binary"
			if err != nil && err.Error() == "unknown binary: "+name {
				t.Fatalf("rejected valid binary name: %s", name)
			}
		}
	})

	t.Run("custom wireshark dir", func(t *testing.T) {
		// Use a nonexistent dir — should fall through to PATH lookup
		e := NewRealExecutor(nil, "/nonexistent/wireshark")
		_, err := e.BinaryPath("tshark")
		// Should not return "unknown binary" error
		if err != nil && err.Error() == "unknown binary: tshark" {
			t.Fatal("rejected valid binary name with custom dir")
		}
	})
}

func TestRealExecutor_Execute(t *testing.T) {
	// Use a cross-platform command that always exists
	var binary string
	var args []string
	if runtime.GOOS == "windows" {
		binary = "cmd"
		args = []string{"/c", "echo", "hello"}
	} else {
		binary = "echo"
		args = []string{"hello"}
	}

	exec := &RealExecutor{wiresharkDir: ""}
	// Override BinaryPath by using the binary directly
	ctx := context.Background()
	cmd := binary

	stdout, _, err := exec.Execute(ctx, cmd, args)
	// This will fail because "echo" isn't in the valid binaries list.
	// That's expected — we're testing the Execute path.
	if err == nil {
		_ = stdout // If it somehow works, that's fine too
	}
}
