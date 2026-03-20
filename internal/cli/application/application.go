package application

import "github.com/spf13/cobra"

// Cmd is the application command group.
var Cmd = &cobra.Command{
	Use:   "application",
	Short: "Application API commands (agent/admin)",
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage authenticated user profile",
}

func init() {
	Cmd.AddCommand(profileCmd)
}
