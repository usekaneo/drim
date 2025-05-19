package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of Kaneo deployment",
	Long:  `Check the status of your Kaneo deployment containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func checkStatus() {
	info := color.New(color.FgCyan).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Checking Kaneo deployment status...\n\n", info("ℹ"))

	cmd := exec.Command("docker", "compose", "ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s Error checking status: %v\n", errorC("✗"), err)
		return
	}
}