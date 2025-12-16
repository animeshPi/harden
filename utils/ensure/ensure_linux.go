package ensure

import (
	"fmt"
	"os"
)

func RequireAdmin() {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Error: This program must be run as root")
		os.Exit(1)
	}
}
