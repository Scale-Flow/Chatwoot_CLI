package accountusers

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List account users",
	RunE:  runList,
}

func init() {
	Cmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
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
	items, err := client.ListAccountUsers(ctx, rctx.AccountID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(items, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
