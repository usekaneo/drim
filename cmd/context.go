package cmd

import (
	"fmt"
	"os"
)

// requireComposeContext fails with drim's own guidance when the working
// directory holds no deployment, instead of letting compose report a missing
// configuration file.
func requireComposeContext() error {
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no Kaneo deployment found in this directory. Run 'drim setup' first, or cd into your deployment directory")
	}
	return nil
}
