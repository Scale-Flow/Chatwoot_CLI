package contacts

import "github.com/spf13/cobra"

// Cmd is the contacts command group.
var Cmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage contacts",
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage contact labels",
}

var contactConversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "List contact conversations",
}

func init() {
	Cmd.AddCommand(labelsCmd)
	Cmd.AddCommand(contactConversationsCmd)
}
