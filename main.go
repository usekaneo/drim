package main

import (
	"github.com/usekaneo/drim/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cmd.SetVersion(Version, BuildTime)
	cmd.Execute()
}
