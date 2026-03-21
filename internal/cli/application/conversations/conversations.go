package conversations

import "github.com/spf13/cobra"

// Cmd is the conversations command group.
var Cmd = &cobra.Command{
	Use:   "conversations",
	Short: "Manage conversations",
}

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage conversation labels",
}

var assignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "Manage conversation assignments",
}

func init() {
	Cmd.AddCommand(labelsCmd)
	Cmd.AddCommand(assignmentsCmd)
}
