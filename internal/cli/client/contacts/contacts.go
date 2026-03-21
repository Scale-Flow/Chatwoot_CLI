package contacts

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "contacts", Short: "Manage client contacts"}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(updateCmd)
}
