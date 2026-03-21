package messages

import "github.com/spf13/cobra"

// Cmd is the messages command group.
var Cmd = &cobra.Command{
	Use:   "messages",
	Short: "Manage messages",
}
