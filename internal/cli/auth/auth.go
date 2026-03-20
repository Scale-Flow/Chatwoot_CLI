package auth

import "github.com/spf13/cobra"

// Cmd is the auth command group.
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication credentials",
}
