package messages

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a message from a conversation",
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	deleteCmd.MarkFlagRequired("conversation-id")
	deleteCmd.Flags().Int("id", 0, "Message ID")
	deleteCmd.MarkFlagRequired("id")
	Cmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
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

	conversationID, _ := cmd.Flags().GetInt("conversation-id")
	id, _ := cmd.Flags().GetInt("id")
	ctx := context.Background()

	if err := client.DeleteMessage(ctx, conversationID, id); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"deleted": true, "id": id})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
