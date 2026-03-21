package inboxes

import "github.com/spf13/cobra"

// Cmd is the inboxes command group.
var Cmd = &cobra.Command{
	Use:   "inboxes",
	Short: "Manage inboxes",
}

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage inbox members",
}

var agentBotCmd = &cobra.Command{
	Use:   "agent-bot",
	Short: "Manage inbox agent bot",
}

func init() {
	Cmd.AddCommand(membersCmd)
	Cmd.AddCommand(agentBotCmd)
}
