package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Kaneo deployment",
	Long:  `Stop the running Kaneo deployment containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopDeployment()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopDeployment() {
	info := color.New(color.FgCyan).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Stopping Kaneo deployment...\n", info("ℹ"))

	cmd := exec.Command("docker", "compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s Error stopping deployment: %v\n", errorC("✗"), err)
		return
	}

	fmt.Printf("%s Kaneo deployment stopped successfully!\n", success("✓"))
}