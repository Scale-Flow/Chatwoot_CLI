package teams

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage teams",
}

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage team members",
}

func init() {
	Cmd.AddCommand(membersCmd)
}
