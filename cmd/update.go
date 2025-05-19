package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Kaneo to the latest version",
	Long: `Update Kaneo containers to the latest version.
This command will pull the latest images and restart the containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateKaneo()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateKaneo() {
	info := color.New(color.FgCyan).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Updating Kaneo to the latest version...\n", info("ℹ"))

	// Pull the latest images
	fmt.Printf("%s Pulling latest images...\n", info("ℹ"))
	pullCmd := exec.Command("docker", "compose", "pull")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	err := pullCmd.Run()
	if err != nil {
		fmt.Printf("%s Error pulling latest images: %v\n", errorC("✗"), err)
		return
	}

	// Restart the containers with the new images
	fmt.Printf("%s Restarting containers with new images...\n", info("ℹ"))
	upCmd := exec.Command("docker", "compose", "up", "-d", "--force-recreate")
	upCmd.Stdout = os.Stdout
	upCmd.Stderr = os.Stderr
	err = upCmd.Run()
	if err != nil {
		fmt.Printf("%s Error restarting containers: %v\n", errorC("✗"), err)
		return
	}

	fmt.Printf("%s Kaneo has been updated to the latest version!\n", success("✓"))
}