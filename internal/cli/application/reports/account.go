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

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Get account report data",
	RunE:  runAccount,
}

func init() {
	accountCmd.Flags().String("metric", "", "Metric name (required)")
	accountCmd.Flags().String("type", "", "Report type (required)")
	accountCmd.Flags().String("since", "", "Start timestamp")
	accountCmd.Flags().String("until", "", "End timestamp")
	accountCmd.Flags().String("id", "", "Resource ID")
	_ = accountCmd.MarkFlagRequired("metric")
	_ = accountCmd.MarkFlagRequired("type")
	Cmd.AddCommand(accountCmd)
}

func runAccount(cmd *cobra.Command, args []string) error {
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

	metric, _ := cmd.Flags().GetString("metric")
	typ, _ := cmd.Flags().GetString("type")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	id, _ := cmd.Flags().GetString("id")

	opts := appapi.ReportOpts{Metric: metric, Type: typ, Since: since, Until: until, ID: id}
	result, err := client.GetReports(context.Background(), opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}
	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
