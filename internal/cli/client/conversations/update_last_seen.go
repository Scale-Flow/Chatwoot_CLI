package conversations

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var updateLastSeenCmd = &cobra.Command{
	Use:   "update-last-seen",
	Short: "Update last seen timestamp in a conversation",
	RunE:  runUpdateLastSeen,
}

func init() {
	updateLastSeenCmd.Flags().Int("conversation-id", 0, "Conversation ID")
	_ = updateLastSeenCmd.MarkFlagRequired("conversation-id")
	Cmd.AddCommand(updateLastSeenCmd)
}

func runUpdateLastSeen(cmd *cobra.Command, args []string) error {
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

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	if err := client.UpdateLastSeen(context.Background(), clientAuth.ContactIdentifier, convID); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"updated": true})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
