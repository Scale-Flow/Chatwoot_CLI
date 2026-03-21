package automationrules

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List automation rules",
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

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	ctx := context.Background()
	rules, err := client.ListAutomationRules(ctx)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(rules, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
