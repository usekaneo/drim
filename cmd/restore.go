package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	backupFile string
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore Kaneo database from backup",
	Long: `Restore the Kaneo SQLite database from a backup file.
This command will restore the database from the specified backup file.`,
	Run: func(cmd *cobra.Command, args []string) {
		restoreDatabase()
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	// Define flags for the restore command
	restoreCmd.Flags().StringVarP(&backupFile, "file", "f", "", "Backup file to restore from")
	restoreCmd.MarkFlagRequired("file")
}

func restoreDatabase() {
	info := color.New(color.FgCyan).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()
	warning := color.New(color.FgYellow).SprintFunc()

	// Check if backup file exists
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		fmt.Printf("%s Backup file not found: %s\n", errorC("✗"), backupFile)
		return
	}

	// Get absolute path to backup file
	absBackupFile, err := filepath.Abs(backupFile)
	if err != nil {
		fmt.Printf("%s Error getting absolute path: %v\n", errorC("✗"), err)
		return
	}

	fmt.Printf("%s This will replace the current database with the backup. All current data will be lost.\n", warning("⚠"))
	fmt.Printf("%s Are you sure you want to continue? (y/n) ", info("?"))

	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
		fmt.Println("Restore cancelled.")
		return
	}

	fmt.Printf("%s Stopping Kaneo containers...\n", info("ℹ"))
	stopCmd := exec.Command("docker", "compose", "stop", "backend")
	err = stopCmd.Run()
	if err != nil {
		fmt.Printf("%s Error stopping containers: %v\n", errorC("✗"), err)
		return
	}

	// Get the container ID
	containerIDCmd := exec.Command("docker", "compose", "ps", "-q", "backend")
	containerIDBytes, err := containerIDCmd.Output()
	if err != nil {
		fmt.Printf("%s Error getting container ID: %v\n", errorC("✗"), err)
		return
	}
	containerID := string(containerIDBytes)
	if containerID == "" {
		fmt.Printf("%s Backend container is not available. Start it first with 'drim deploy'.\n", errorC("✗"))
		return
	}

	// Copy the backup file to the container
	fmt.Printf("%s Restoring database from backup...\n", info("ℹ"))
	copyCmd := exec.Command("docker", "cp",
		absBackupFile,
		fmt.Sprintf("%s:/app/apps/api/data/kaneo.db", containerID[:12]))
	err = copyCmd.Run()
	if err != nil {
		fmt.Printf("%s Error restoring database: %v\n", errorC("✗"), err)
		return
	}

	// Start the containers again
	fmt.Printf("%s Starting Kaneo containers...\n", info("ℹ"))
	startCmd := exec.Command("docker", "compose", "start", "backend")
	err = startCmd.Run()
	if err != nil {
		fmt.Printf("%s Error starting containers: %v\n", errorC("✗"), err)
		return
	}

	fmt.Printf("%s Database restored successfully!\n", success("✓"))
}