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

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a platform account",
	RunE:  runGet,
}

func init() {
	Cmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
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
	account, err := client.GetAccount(ctx, rctx.AccountID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(account)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
