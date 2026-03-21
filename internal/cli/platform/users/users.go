package users

import "github.com/spf13/cobra"

// Cmd is the users command group.
var Cmd = &cobra.Command{
	Use:   "users",
	Short: "Manage platform users",
}
