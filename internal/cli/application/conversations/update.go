package conversations

import (
	"context"
	"fmt"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a conversation",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Conversation ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("status", "", "Conversation status")
	updateCmd.Flags().String("priority", "", "Conversation priority")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	statusChanged := cmd.Flags().Changed("status")
	priorityChanged := cmd.Flags().Changed("priority")
	if !statusChanged && !priorityChanged {
		return fmt.Errorf("requires at least one of --status or --priority")
	}

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
	opts := appapi.UpdateConversationOpts{}

	if statusChanged {
		v, _ := cmd.Flags().GetString("status")
		opts.Status = &v
	}
	if priorityChanged {
		v, _ := cmd.Flags().GetString("priority")
		opts.Priority = &v
	}

	ctx := context.Background()
	convo, err := client.UpdateConversation(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(convo)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
