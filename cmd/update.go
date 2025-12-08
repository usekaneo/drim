package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/ui"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update drim to the latest version",
	Long:  `Downloads and installs the latest version of drim from GitHub releases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Checking for updates...")

		currentPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		currentPath, err = filepath.EvalSymlinks(currentPath)
		if err != nil {
			return fmt.Errorf("failed to resolve executable path: %w", err)
		}

		binaryName := getBinaryName()
		downloadURL := fmt.Sprintf("https://github.com/usekaneo/drim/releases/latest/download/%s", binaryName)

		ui.Info(fmt.Sprintf("Downloading latest version from %s...", downloadURL))

		resp, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download update: HTTP %d", resp.StatusCode)
		}

		tmpFile, err := os.CreateTemp("", "drim-update-*")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to save update: %w", err)
		}
		tmpFile.Close()

		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			return fmt.Errorf("failed to make update executable: %w", err)
		}

		if err := os.Rename(tmpFile.Name(), currentPath); err != nil {
			return fmt.Errorf("failed to replace executable: %w (you may need to run with sudo)", err)
		}

		ui.Success("✨ drim updated successfully!")
		ui.Info("Run 'drim --version' to verify the new version")
		return nil
	},
}

func getBinaryName() string {
	var platform, arch string

	switch runtime.GOOS {
	case "linux":
		platform = "Linux"
	case "darwin":
		platform = "Darwin"
	case "windows":
		platform = "Windows"
	default:
		platform = "Linux"
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "arm64"
	default:
		arch = "x86_64"
	}

	if platform == "Windows" {
		return fmt.Sprintf("drim_%s_%s.exe", platform, arch)
	}
	return fmt.Sprintf("drim_%s_%s", platform, arch)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

