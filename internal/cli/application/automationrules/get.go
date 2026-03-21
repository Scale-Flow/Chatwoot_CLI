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

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an automation rule",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().Int("id", 0, "Automation rule ID")
	getCmd.MarkFlagRequired("id")
	Cmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
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

	id, _ := cmd.Flags().GetInt("id")
	ctx := context.Background()

	rule, err := client.GetAutomationRule(ctx, id)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(rule)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
