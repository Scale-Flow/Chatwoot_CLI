package contacts

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge two contacts",
	RunE:  runMerge,
}

func init() {
	mergeCmd.Flags().Int("base-id", 0, "Base contact ID (kept)")
	mergeCmd.MarkFlagRequired("base-id")
	mergeCmd.Flags().Int("merge-id", 0, "Contact ID to merge (discarded)")
	mergeCmd.MarkFlagRequired("merge-id")
	Cmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	baseID, _ := cmd.Flags().GetInt("base-id")
	mergeID, _ := cmd.Flags().GetInt("merge-id")
	ctx := context.Background()

	contact, err := client.MergeContacts(ctx, baseID, mergeID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(contact)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
