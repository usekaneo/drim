package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drim",
	Short: "Drim is a CLI tool for deploying Kaneo",
	Long: `Drim is a command-line interface tool that helps you deploy Kaneo,
an open source project management platform focused on simplicity and efficiency.
It simplifies the deployment process using Docker and Docker Compose.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "completion" && cmd.Name() != "__complete" {
			printLogo()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printLogo() {
	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	logo := `
    ____       _
   / __ \_____(_)___ ___
  / / / / ___/ / __ ` + "`" + `__ \
 / /_/ / /  / / / / / / /
/_____/_/  /_/_/ /_/ /_/
                           `

	fmt.Println(cyan(logo))
	fmt.Printf("%s CLI tool for %s\n\n", blue("Drim"), yellow("Kaneo"))
}
