package banner

import (
	"fmt"
)

const (
	ColorReset = "\033[0m"
	ColorCyan  = "\033[36m"
	ColorBold  = "\033[1m"
)

func Print() {
	banner := `
    ____  ____  ________  ___
   / __ \/ __ \/  _/ __ \/   |
  / / / / /_/ // // / / / /| |
 / /_/ / _, _// // /_/ / ___ |
/_____/_/ |_/___/_____/_/  |_|

One-Click Self-Hosted Kaneo Deployment
`
	fmt.Printf("%s%s%s%s\n", ColorBold, ColorCyan, banner, ColorReset)
}
