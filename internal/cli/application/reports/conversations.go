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

var conversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "Get conversation metrics",
	RunE:  runConversations,
}

func init() {
	conversationsCmd.Flags().String("type", "", "Report type: account or agent (required)")
	conversationsCmd.Flags().String("since", "", "Start timestamp")
	conversationsCmd.Flags().String("until", "", "End timestamp")
	_ = conversationsCmd.MarkFlagRequired("type")
	Cmd.AddCommand(conversationsCmd)
}

func runConversations(cmd *cobra.Command, args []string) error {
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

	typ, _ := cmd.Flags().GetString("type")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")

	opts := appapi.ReportOpts{Type: typ, Since: since, Until: until}
	ctx := context.Background()
	var result any
	if typ == "agent" {
		result, err = client.GetAgentConversationMetrics(ctx, opts)
	} else {
		result, err = client.GetConversationMetrics(ctx, opts)
	}
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}
	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
