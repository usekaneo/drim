package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/banner"
)

var (
	version   string
	buildTime string
)

var rootCmd = &cobra.Command{
	Use:   "drim",
	Short: "One-Click Self-Hosted Kaneo Deployment Tool",
	Long: `drim is a simple, fast command-line tool that lets anyone deploy 
and manage a full Kaneo stack with a single command.

It automatically sets up PostgreSQL, Kaneo API, Kaneo Web, and Caddy 
as a reverse proxy with automatic HTTPS.`,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
		cmd.Help()
	},
}

func SetVersion(v, bt string) {
	version = v
	buildTime = bt
	rootCmd.Version = fmt.Sprintf("%s (built: %s)", v, bt)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
