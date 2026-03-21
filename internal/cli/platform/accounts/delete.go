package accounts

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a platform account",
	RunE:  runDelete,
}

func init() {
	Cmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}
	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModePlatform)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}
	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := platapi.NewClient(transport)

	ctx := context.Background()
	if err := client.DeleteAccount(ctx, rctx.AccountID); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"deleted": true, "id": rctx.AccountID})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
