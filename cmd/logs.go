package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	follow bool
	tail   string
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View logs from Kaneo containers",
	Long: `View logs from Kaneo containers.
You can specify a service name (backend or frontend) or view all logs.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viewLogs(args)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVarP(&tail, "tail", "t", "all", "Number of lines to show from the end of the logs")
}

func viewLogs(args []string) {
	info := color.New(color.FgCyan).SprintFunc()
	errorC := color.New(color.FgRed).SprintFunc()

	service := ""
	if len(args) > 0 {
		service = args[0]
		if service != "backend" && service != "frontend" && service != "traefik" {
			fmt.Printf("%s Invalid service name. Use 'backend', 'frontend', or 'traefik'.\n", errorC("✗"))
			return
		}
	}

	fmt.Printf("%s Viewing logs for %s...\n\n", info("ℹ"), service)

	dockerArgs := []string{"compose", "logs"}

	if follow {
		dockerArgs = append(dockerArgs, "--follow")
	}

	if tail != "all" {
		dockerArgs = append(dockerArgs, "--tail", tail)
	}

	if service != "" {
		dockerArgs = append(dockerArgs, service)
	}

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n%s Error viewing logs: %v\n", errorC("✗"), err)
		return
	}
}
