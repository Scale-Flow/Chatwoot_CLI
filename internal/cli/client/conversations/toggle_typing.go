package conversations

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var toggleTypingCmd = &cobra.Command{
	Use:   "toggle-typing",
	Short: "Toggle typing status in a conversation",
	RunE:  runToggleTyping,
}

func init() {
	toggleTypingCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	_ = toggleTypingCmd.MarkFlagRequired("conversation-id")
	toggleTypingCmd.Flags().String("status", "", "Typing status: on or off")
	_ = toggleTypingCmd.MarkFlagRequired("status")
	Cmd.AddCommand(toggleTypingCmd)
}

func runToggleTyping(cmd *cobra.Command, args []string) error {
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
	status, _ := cmd.Flags().GetString("status")

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	opts := clientapi.ToggleTypingOpts{TypingStatus: status}
	if err := client.ToggleTyping(context.Background(), clientAuth.ContactIdentifier, convID, opts); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"typing_status": status})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
