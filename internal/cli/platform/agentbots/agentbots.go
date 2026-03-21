package agentbots

import "github.com/spf13/cobra"

// Cmd is the agent-bots command group.
var Cmd = &cobra.Command{
	Use:   "agent-bots",
	Short: "Manage platform agent bots",
}
