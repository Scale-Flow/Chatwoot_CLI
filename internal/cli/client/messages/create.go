package messages

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a message in a conversation",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	_ = createCmd.MarkFlagRequired("conversation-id")
	createCmd.Flags().String("content", "", "Message content")
	_ = createCmd.MarkFlagRequired("content")
	createCmd.Flags().String("type", "", "Message type")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
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
	content, _ := cmd.Flags().GetString("content")
	msgType, _ := cmd.Flags().GetString("type")

	opts := clientapi.CreateMessageOpts{
		Content:     content,
		MessageType: msgType,
	}

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	message, err := client.CreateMessage(context.Background(), clientAuth.ContactIdentifier, convID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(message)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
