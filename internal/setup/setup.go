package setup

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	wiresharkDownloadPage = "https://www.wireshark.org/download.html"
	wiresharkMirror       = "https://2.na.dl.wireshark.org"
	npcapURL              = "https://npcap.com/dist/npcap-1.80.exe"
)

// Status describes the current tshark installation state.
type Status struct {
	Installed      bool
	TSharkPath     string
	Version        string
	CapinfosPath   string
	EditcapPath    string
	MergecapPath   string
	NpcapInstalled bool // Windows only: whether Npcap is installed for live capture
}

// Check detects whether tshark and related tools are installed.
func Check(logger *slog.Logger) *Status {
	s := &Status{}

	binaries := map[string]*string{
		"tshark":   &s.TSharkPath,
		"capinfos": &s.CapinfosPath,
		"editcap":  &s.EditcapPath,
		"mergecap": &s.MergecapPath,
	}

	for name, pathPtr := range binaries {
		path, err := findBinary(name)
		if err == nil {
			*pathPtr = path
		}
	}

	s.Installed = s.TSharkPath != ""

	if s.Installed {
		if version, err := getTSharkVersion(s.TSharkPath); err == nil {
			s.Version = version
		}
	}

	if runtime.GOOS == "windows" {
		s.NpcapInstalled = checkNpcap()
	}

	return s
}

// Install attempts to install Wireshark/tshark for the current platform.
func Install(ctx context.Context, logger *slog.Logger) error {
	switch runtime.GOOS {
	case "windows":
		return installWindows(ctx, logger)
	case "darwin":
		return installDarwin(ctx, logger)
	case "linux":
		return installLinux(ctx, logger)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// PrintInstructions prints manual installation instructions for the current platform.
func PrintInstructions() string {
	var sb strings.Builder

	sb.WriteString("Wireshark/tshark Installation Guide\n")
	sb.WriteString("===================================\n\n")
	sb.WriteString("Reference: https://tshark.dev/setup/install/\n\n")

	switch runtime.GOOS {
	case "windows":
		sb.WriteString("Windows:\n")
		sb.WriteString("  1. Download from https://www.wireshark.org/download.html\n")
		sb.WriteString("  2. Run installer (or use: sharkline setup --install)\n")
		sb.WriteString("  3. Ensure 'Add to PATH' is checked during installation\n")
		sb.WriteString("  4. Npcap will be installed automatically for live capture\n\n")
		sb.WriteString("  Silent install:\n")
		sb.WriteString("    Wireshark-win64-latest.exe /S /D=C:\\Program Files\\Wireshark\n\n")
		sb.WriteString("  Chocolatey:\n")
		sb.WriteString("    choco install wireshark\n\n")
		sb.WriteString("  Winget:\n")
		sb.WriteString("    winget install WiresharkFoundation.Wireshark\n")

	case "darwin":
		sb.WriteString("macOS:\n")
		sb.WriteString("  Homebrew (recommended):\n")
		sb.WriteString("    brew install --cask wireshark\n\n")
		sb.WriteString("  Or download from https://www.wireshark.org/download.html\n")

	case "linux":
		sb.WriteString("Linux:\n")
		sb.WriteString("  Debian/Ubuntu:\n")
		sb.WriteString("    sudo apt-get update && sudo apt-get install -y tshark\n\n")
		sb.WriteString("  Fedora/CentOS/RHEL:\n")
		sb.WriteString("    sudo dnf install wireshark-cli\n\n")
		sb.WriteString("  Arch:\n")
		sb.WriteString("    sudo pacman -S wireshark-cli\n\n")
		sb.WriteString("  Alpine:\n")
		sb.WriteString("    sudo apk add tshark\n")

	default:
		sb.WriteString("See https://tshark.dev/setup/install/ for your platform.\n")
	}

	return sb.String()
}

func installWindows(ctx context.Context, logger *slog.Logger) error {
	// Try winget first (most modern Windows approach).
	if wingetPath, err := exec.LookPath("winget"); err == nil {
		logger.Info("installing via winget", "path", wingetPath)
		cmd := exec.CommandContext(ctx, "winget", "install", "--id", "WiresharkFoundation.Wireshark", "--accept-package-agreements", "--accept-source-agreements")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
		logger.Warn("winget install failed, trying direct download")
	}

	// Try chocolatey.
	if chocoPath, err := exec.LookPath("choco"); err == nil {
		logger.Info("installing via chocolatey", "path", chocoPath)
		cmd := exec.CommandContext(ctx, "choco", "install", "wireshark", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
		logger.Warn("chocolatey install failed, trying direct download")
	}

	// Direct download fallback.
	logger.Info("downloading Wireshark installer")

	version, err := fetchLatestVersion(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to detect latest version: %w", err)
	}

	url := fmt.Sprintf("%s/win64/Wireshark-4.6.%s-x64.exe", wiresharkMirror, version)
	installerPath := filepath.Join(os.TempDir(), "wireshark-installer.exe")

	if err := downloadFile(ctx, url, installerPath, logger); err != nil {
		// Fallback to a known good version URL.
		url = fmt.Sprintf("%s/win64/Wireshark-4.6.4-x64.exe", wiresharkMirror)
		if err := downloadFile(ctx, url, installerPath, logger); err != nil {
			return fmt.Errorf("download failed: %w\n\nManual install:\n%s", err, PrintInstructions())
		}
	}

	logger.Info("running silent installer", "path", installerPath)
	cmd := exec.CommandContext(ctx, installerPath, "/S", "/desktopicon=no")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installer failed: %w\n\nRun manually: %s /S", err, installerPath)
	}

	logger.Info("Wireshark installed successfully")
	return nil
}

func installDarwin(ctx context.Context, logger *slog.Logger) error {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("Homebrew not found. Install it from https://brew.sh then run:\n  brew install --cask wireshark")
	}

	logger.Info("installing via Homebrew", "path", brewPath)
	cmd := exec.CommandContext(ctx, "brew", "install", "--cask", "wireshark")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installLinux(ctx context.Context, logger *slog.Logger) error {
	// Detect package manager.
	type pkgMgr struct {
		binary string
		args   []string
	}

	managers := []pkgMgr{
		{"apt-get", []string{"install", "-y", "tshark"}},
		{"dnf", []string{"install", "-y", "wireshark-cli"}},
		{"yum", []string{"install", "-y", "wireshark-cli"}},
		{"pacman", []string{"-S", "--noconfirm", "wireshark-cli"}},
		{"apk", []string{"add", "tshark"}},
	}

	for _, mgr := range managers {
		if _, err := exec.LookPath(mgr.binary); err != nil {
			continue
		}

		logger.Info("installing via package manager", "manager", mgr.binary)

		// Update package lists for apt.
		if mgr.binary == "apt-get" {
			update := exec.CommandContext(ctx, "sudo", "apt-get", "update")
			update.Stdout = os.Stdout
			update.Stderr = os.Stderr
			_ = update.Run()
		}

		args := append([]string{mgr.binary}, mgr.args...)
		cmd := exec.CommandContext(ctx, "sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return fmt.Errorf("no supported package manager found\n\n%s", PrintInstructions())
}

func findBinary(name string) (string, error) {
	// Try PATH first.
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}

	// Windows-specific locations.
	if runtime.GOOS == "windows" {
		suffix := ".exe"
		for _, dir := range []string{
			`C:\Program Files\Wireshark`,
			`C:\Program Files (x86)\Wireshark`,
		} {
			candidate := filepath.Join(dir, name+suffix)
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	// macOS Homebrew location.
	if runtime.GOOS == "darwin" {
		for _, dir := range []string{
			"/usr/local/bin",
			"/opt/homebrew/bin",
			"/Applications/Wireshark.app/Contents/MacOS",
		} {
			candidate := filepath.Join(dir, name)
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("%s not found", name)
}

func getTSharkVersion(tsharkPath string) (string, error) {
	out, err := exec.Command(tsharkPath, "--version").Output()
	if err != nil {
		return "", err
	}
	// First line: "TShark (Wireshark) X.Y.Z ..."
	line := strings.Split(string(out), "\n")[0]
	re := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	if m := re.FindString(line); m != "" {
		return m, nil
	}
	return strings.TrimSpace(line), nil
}

func fetchLatestVersion(ctx context.Context, logger *slog.Logger) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", wiresharkDownloadPage, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil {
		return "", err
	}

	// Extract version from download links like "Wireshark-4.6.4-x64.exe"
	re := regexp.MustCompile(`Wireshark-(\d+\.\d+\.\d+)-x64\.exe`)
	if m := re.FindSubmatch(body); len(m) > 1 {
		version := string(m[1])
		logger.Info("detected latest Wireshark version", "version", version)
		return version, nil
	}

	return "", fmt.Errorf("could not detect version from download page")
}

func downloadFile(ctx context.Context, url string, dest string, logger *slog.Logger) error {
	logger.Info("downloading", "url", url, "dest", dest)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	written, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	logger.Info("download complete", "bytes", written)
	return nil
}
