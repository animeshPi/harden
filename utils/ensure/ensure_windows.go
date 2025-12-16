package ensure

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func RequireAdmin() {
	var token windows.Token

	err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_QUERY,
		&token,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to query process token")
		os.Exit(1)
	}
	defer token.Close()

	adminSID, err := windows.CreateWellKnownSid(
		windows.WinBuiltinAdministratorsSid,
		nil,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to create admin SID")
		os.Exit(1)
	}

	ok, err := token.IsMember(adminSID)
	if err != nil || !ok {
		fmt.Fprintln(os.Stderr, "Error: This program must be run as Administrator")
		os.Exit(1)
	}
}
