package messages

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "messages", Short: "Manage client messages"}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
}
