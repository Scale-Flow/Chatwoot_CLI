package customfilters

import "github.com/spf13/cobra"

func init() {
	root := &cobra.Command{Use: "chatwoot", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().Bool("pretty", false, "Indent JSON output")
	root.PersistentFlags().String("profile", "", "Select named profile")
	root.PersistentFlags().String("base-url", "", "Override base URL")
	root.PersistentFlags().Int("account-id", 0, "Override account ID")
	appCmd := &cobra.Command{Use: "application"}
	appCmd.AddCommand(Cmd)
	root.AddCommand(appCmd)
}
