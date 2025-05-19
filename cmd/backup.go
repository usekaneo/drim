package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	backupDir string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup Kaneo database",
	Long: `Backup the Kaneo SQLite database.
This command will create a backup of the database in the specified directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		backupDatabase()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().StringVarP(&backupDir, "dir", "d", "./backups", "Directory to store backups")
}

func backupDatabase() {
	info := color.New(color.FgCyan).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			fmt.Printf("%s Error creating backup directory: %v\n", errorC("✗"), err)
			return
		}
	}

	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("kaneo-backup-%s.db", timestamp))

	fmt.Printf("%s Creating backup of Kaneo database...\n", info("ℹ"))

	containerIDCmd := exec.Command("docker", "compose", "ps", "-q", "backend")
	containerIDBytes, err := containerIDCmd.Output()
	if err != nil {
		fmt.Printf("%s Error getting container ID: %v\n", errorC("✗"), err)
		return
	}
	containerID := string(containerIDBytes)
	if containerID == "" {
		fmt.Printf("%s Backend container is not running. Start it first with 'drim deploy'.\n", errorC("✗"))
		return
	}

	copyCmd := exec.Command("docker", "cp",
		fmt.Sprintf("%s:/app/apps/api/data/kaneo.db", containerID[:12]),
		backupFile)
	err = copyCmd.Run()
	if err != nil {
		fmt.Printf("%s Error copying database file: %v\n", errorC("✗"), err)
		return
	}

	fmt.Printf("%s Backup created successfully: %s\n", success("✓"), backupFile)
	fmt.Printf("%s To restore this backup, replace the database file in the container or volume.\n", info("ℹ"))
}
