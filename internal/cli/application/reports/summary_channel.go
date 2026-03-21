package reports

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var summaryChannelCmd = &cobra.Command{
	Use:   "summary-by-channel",
	Short: "Get summary report grouped by channel",
	RunE:  runSummaryChannel,
}

func init() {
	summaryChannelCmd.Flags().String("since", "", "Start timestamp")
	summaryChannelCmd.Flags().String("until", "", "End timestamp")
	Cmd.AddCommand(summaryChannelCmd)
}

func runSummaryChannel(cmd *cobra.Command, args []string) error {
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

	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")

	opts := appapi.ReportOpts{Since: since, Until: until}
	result, err := client.GetSummaryByChannel(context.Background(), opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}
	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
