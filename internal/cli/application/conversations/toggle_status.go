package conversations

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var toggleStatusCmd = &cobra.Command{
	Use:   "toggle-status",
	Short: "Toggle conversation status",
	RunE:  runToggleStatus,
}

func init() {
	toggleStatusCmd.Flags().Int("id", 0, "Conversation ID")
	toggleStatusCmd.MarkFlagRequired("id")
	toggleStatusCmd.Flags().String("status", "", "Status (open, resolved, pending, snoozed)")
	toggleStatusCmd.MarkFlagRequired("status")
	Cmd.AddCommand(toggleStatusCmd)
}

func runToggleStatus(cmd *cobra.Command, args []string) error {
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
	status, _ := cmd.Flags().GetString("status")
	ctx := context.Background()

	convo, err := client.ToggleConversationStatus(ctx, id, status)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(convo)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
