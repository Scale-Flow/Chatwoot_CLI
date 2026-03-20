package auth

import (
	"github.com/spf13/cobra"
)

func init() {
	root := &cobra.Command{
		Use:           "chatwoot",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().Bool("pretty", false, "Indent JSON output")
	root.PersistentFlags().String("profile", "", "Select named profile")
	root.AddCommand(Cmd)
}
