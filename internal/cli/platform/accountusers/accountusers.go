package accountusers

import "github.com/spf13/cobra"

// Cmd is the account-users command group.
var Cmd = &cobra.Command{
	Use:   "account-users",
	Short: "Manage account users",
}
