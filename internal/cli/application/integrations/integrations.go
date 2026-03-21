package integrations

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "integrations", Short: "Manage integrations"}

var appsCmd = &cobra.Command{Use: "apps", Short: "Manage integration apps"}

var hooksCmd = &cobra.Command{Use: "hooks", Short: "Manage integration hooks"}

func init() {
	Cmd.AddCommand(appsCmd)
	Cmd.AddCommand(hooksCmd)
}
