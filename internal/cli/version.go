package cli

import (
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build and version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp := contract.Success(version.Info())
		return contract.Write(cmd.OutOrStdout(), resp, prettyFlag)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
