package executor

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// CommandExecutor is the primary port for executing Wireshark CLI tools.
// All tool handlers depend on this interface, enabling testing without
// real tshark binaries installed.
type CommandExecutor interface {
	// Execute runs a CLI command and returns stdout, stderr, and any error.
	Execute(ctx context.Context, binary string, args []string) (stdout []byte, stderr []byte, err error)

	// BinaryPath resolves the full path to a Wireshark CLI tool.
	// Supports: tshark, capinfos, editcap, mergecap, dumpcap.
	BinaryPath(name string) (string, error)
}

// DefaultTimeout is the maximum execution time for any single command.
const DefaultTimeout = 60 * time.Second

// RealExecutor executes commands against real Wireshark CLI binaries.
type RealExecutor struct {
	logger *slog.Logger
	// wiresharkDir overrides the default Wireshark install path.
	// If empty, uses PATH lookup and platform-specific defaults.
	wiresharkDir string
}

// NewRealExecutor creates a new executor that runs real CLI commands.
func NewRealExecutor(logger *slog.Logger, wiresharkDir string) *RealExecutor {
	return &RealExecutor{
		logger:       logger,
		wiresharkDir: wiresharkDir,
	}
}

func (e *RealExecutor) Execute(ctx context.Context, binary string, args []string) ([]byte, []byte, error) {
	path, err := e.BinaryPath(binary)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}

	e.logger.Debug("executing command", "binary", path, "args", args)

	cmd := exec.CommandContext(ctx, path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return stdout.Bytes(), stderr.Bytes(), fmt.Errorf("%s failed: %w (stderr: %s)", binary, err, strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}

func (e *RealExecutor) BinaryPath(name string) (string, error) {
	validBinaries := map[string]bool{
		"tshark":   true,
		"capinfos": true,
		"editcap":  true,
		"mergecap": true,
		"dumpcap":  true,
	}

	if !validBinaries[name] {
		return "", fmt.Errorf("unknown binary: %s", name)
	}

	// If a custom Wireshark dir is set, look there first.
	if e.wiresharkDir != "" {
		suffix := ""
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}
		candidate := e.wiresharkDir + "/" + name + suffix
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}

	// Try PATH lookup.
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	// Platform-specific defaults.
	if runtime.GOOS == "windows" {
		for _, dir := range []string{
			`C:\Program Files\Wireshark`,
			`C:\Program Files (x86)\Wireshark`,
		} {
			candidate := dir + `\` + name + ".exe"
			if _, err := exec.LookPath(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("%s not found: ensure Wireshark is installed and in PATH", name)
}
