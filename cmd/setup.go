package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/inovacc/sharkline/internal/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Check and install Wireshark/tshark dependencies",
	Long: `Check whether tshark and related Wireshark CLI tools are installed,
and optionally install them automatically.

Without flags, shows the current installation status.
Use --install to auto-install for your platform.

Supported installation methods:
  Windows:  winget → chocolatey → direct download (silent install)
  macOS:    Homebrew (brew install --cask wireshark)
  Linux:    apt-get → dnf → yum → pacman → apk

Reference: https://tshark.dev/setup/install/`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		install, _ := cmd.Flags().GetBool("install")

		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		// Check current status.
		status := setup.Check(logger)

		fmt.Println("Sharkline Dependency Check")
		fmt.Println("==========================")
		fmt.Printf("Platform: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

		type tool struct {
			name string
			path string
		}
		tools := []tool{
			{"tshark", status.TSharkPath},
			{"capinfos", status.CapinfosPath},
			{"editcap", status.EditcapPath},
			{"mergecap", status.MergecapPath},
		}

		allFound := true
		for _, t := range tools {
			if t.path != "" {
				fmt.Printf("  [OK]      %-10s %s\n", t.name, t.path)
			} else {
				fmt.Printf("  [MISSING] %-10s not found\n", t.name)
				allFound = false
			}
		}

		if status.Version != "" {
			fmt.Printf("\n  Version: %s\n", status.Version)
		}

		if allFound {
			fmt.Println("\nAll dependencies are installed. Sharkline is ready to use.")
			return nil
		}

		fmt.Println()

		if !install {
			fmt.Println(setup.PrintInstructions())
			fmt.Println("Or run: sharkline setup --install")
			return nil
		}

		fmt.Println("Installing Wireshark/tshark...")
		if err := setup.Install(cmd.Context(), logger); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		// Re-check after install.
		fmt.Println("\nVerifying installation...")
		postStatus := setup.Check(logger)
		if postStatus.Installed {
			fmt.Printf("tshark %s installed at %s\n", postStatus.Version, postStatus.TSharkPath)
			fmt.Println("Sharkline is ready to use.")
		} else {
			fmt.Println("tshark not detected after install. You may need to restart your terminal or add Wireshark to your PATH.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().Bool("install", false, "auto-install Wireshark/tshark for your platform")
}
