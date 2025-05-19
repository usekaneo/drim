package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showcaseCmd = &cobra.Command{
	Use:   "showcase",
	Short: "Show a cool demo of Kaneo deployment",
	Long:  `Create a visually appealing showcase of Kaneo deployment for social media sharing.`,
	Run: func(cmd *cobra.Command, args []string) {
		runShowcase()
	},
}

func init() {
	rootCmd.AddCommand(showcaseCmd)
}

func runShowcase() {
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	clearScreen()

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	features := []string{
		"🚀 Easy deployment with multiple proxy options",
		"🔒 Automatic HTTPS configuration",
		"💾 Database backup and restore",
		"📊 Container status monitoring",
		"📝 Real-time log viewing",
		"🔄 One-click updates",
	}

	logo := `
    ____       _
   / __ \_____(_)___ ___
  / / / / ___/ / __ ` + "`" + `__ \
 / /_/ / /  / / / / / / /
/_____/_/  /_/_/ /_/ /_/
                           `

	fmt.Printf("\n%s\n", cyan(logo))
	fmt.Printf("%s %s %s\n\n", magenta("●"), yellow("●"), green("●"))

	for _, feature := range features {
		for i := 0; i < 5; i++ {
			fmt.Printf("\r%s %s", frames[i%len(frames)], feature)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("\r✓ %s\n", bold(feature))
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Printf("\n%s Deploy your Kaneo instance in seconds:\n\n", yellow("→"))
	fmt.Printf("   %s\n", cyan("$ drim deploy"))
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("   %s\n", magenta("$ drim deploy --domain example.com --proxy traefik --https"))
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("\n%s Monitor and manage with ease:\n\n", yellow("→"))
	fmt.Printf("   %s\n", cyan("$ drim status"))
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   %s\n", cyan("$ drim logs"))
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   %s\n", cyan("$ drim backup"))

	fmt.Printf("\n%s Start your journey: %s\n\n", green("→"), bold("github.com/usekaneo/kaneo"))
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
