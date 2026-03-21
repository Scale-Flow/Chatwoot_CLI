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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a message in a conversation",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	createCmd.MarkFlagRequired("conversation-id")
	createCmd.Flags().String("content", "", "Message content")
	createCmd.MarkFlagRequired("content")
	createCmd.Flags().String("message-type", "outgoing", "Message type (incoming, outgoing)")
	createCmd.Flags().Bool("private", false, "Send as private note")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
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
	content, _ := cmd.Flags().GetString("content")
	msgType, _ := cmd.Flags().GetString("message-type")
	private, _ := cmd.Flags().GetBool("private")

	opts := appapi.CreateMessageOpts{
		Content:     content,
		MessageType: msgType,
		Private:     private,
	}

	ctx := context.Background()
	msg, err := client.CreateMessage(ctx, conversationID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(msg)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
