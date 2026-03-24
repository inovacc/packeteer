package setup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	npcapVersion     = "1.80"
	npcapDownloadURL = "https://npcap.com/dist/npcap-" + npcapVersion + ".exe"
)

// checkNpcap returns true if Npcap is installed on Windows.
func checkNpcap() bool {
	if runtime.GOOS != "windows" {
		return true // Not applicable on non-Windows.
	}

	// Check for Npcap DLL in System32.
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}

	npcapDLL := filepath.Join(systemRoot, "System32", "Npcap", "wpcap.dll")
	if _, err := os.Stat(npcapDLL); err == nil {
		return true
	}

	// Fallback: check the standard Npcap install directory.
	npcapDir := filepath.Join(systemRoot, "System32", "Npcap")
	if _, err := os.Stat(npcapDir); err == nil {
		return true
	}

	return false
}

// InstallNpcap downloads and installs Npcap on Windows.
func InstallNpcap(ctx context.Context, logger *slog.Logger) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Npcap is only needed on Windows")
	}

	if checkNpcap() {
		logger.Info("Npcap is already installed")
		return nil
	}

	// Try winget first.
	if _, err := exec.LookPath("winget"); err == nil {
		logger.Info("installing Npcap via winget")
		cmd := exec.CommandContext(ctx, "winget", "install", "--id", "Insecure.Npcap",
			"--accept-package-agreements", "--accept-source-agreements")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
		logger.Warn("winget install failed, trying direct download")
	}

	// Direct download.
	installerPath := filepath.Join(os.TempDir(), "npcap-installer.exe")
	if err := downloadFile(ctx, npcapDownloadURL, installerPath, logger); err != nil {
		return fmt.Errorf("failed to download Npcap: %w\n\nDownload manually from https://npcap.com/", err)
	}

	logger.Info("running Npcap installer", "path", installerPath)
	cmd := exec.CommandContext(ctx, installerPath, "/S")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Npcap installer failed: %w\n\nRun manually: %s", err, installerPath)
	}

	logger.Info("Npcap installed successfully")
	return nil
}

// NpcapInstructions returns manual install instructions for Npcap.
func NpcapInstructions() string {
	return `Npcap Installation (required for live packet capture on Windows)
================================================================

Npcap provides the packet capture driver that tshark needs for live capture.
Without it, you can still read/analyze existing pcap files.

  Download: https://npcap.com/
  Winget:   winget install Insecure.Npcap
  Auto:     sharkline setup --install-npcap
`
}
