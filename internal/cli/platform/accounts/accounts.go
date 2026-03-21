package accounts

import "github.com/spf13/cobra"

// Cmd is the accounts command group.
var Cmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage platform accounts",
}
