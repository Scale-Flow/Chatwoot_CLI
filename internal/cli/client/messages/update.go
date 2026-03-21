package messages

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a message in a conversation",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	_ = updateCmd.MarkFlagRequired("conversation-id")
	updateCmd.Flags().Int("message-id", 0, "Message ID")
	_ = updateCmd.MarkFlagRequired("message-id")
	updateCmd.Flags().String("content", "", "Message content")
	_ = updateCmd.MarkFlagRequired("content")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}
	clientAuth, err := cmdutil.ResolveClientAuth(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}
	if clientAuth.ContactIdentifier == "" {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, "--contact-id flag or CHATWOOT_CONTACT_IDENTIFIER env var required for this command")
	}

	convID, _ := cmd.Flags().GetInt("conversation-id")
	msgID, _ := cmd.Flags().GetInt("message-id")
	content, _ := cmd.Flags().GetString("content")

	opts := clientapi.UpdateMessageOpts{Content: content}

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	message, err := client.UpdateMessage(context.Background(), clientAuth.ContactIdentifier, convID, msgID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(message)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
