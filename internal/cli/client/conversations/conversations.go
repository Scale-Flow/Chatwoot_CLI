package conversations

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "conversations", Short: "Manage client conversations"}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(toggleStatusCmd)
	Cmd.AddCommand(toggleTypingCmd)
	Cmd.AddCommand(updateLastSeenCmd)
}
