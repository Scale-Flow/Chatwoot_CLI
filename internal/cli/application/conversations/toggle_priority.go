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

var togglePriorityCmd = &cobra.Command{
	Use:   "toggle-priority",
	Short: "Toggle conversation priority",
	RunE:  runTogglePriority,
}

func init() {
	togglePriorityCmd.Flags().Int("id", 0, "Conversation ID")
	togglePriorityCmd.MarkFlagRequired("id")
	togglePriorityCmd.Flags().String("priority", "", "Priority (urgent, high, medium, low, none)")
	togglePriorityCmd.MarkFlagRequired("priority")
	Cmd.AddCommand(togglePriorityCmd)
}

func runTogglePriority(cmd *cobra.Command, args []string) error {
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
	priority, _ := cmd.Flags().GetString("priority")
	ctx := context.Background()

	result, err := client.ToggleConversationPriority(ctx, id, priority)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
