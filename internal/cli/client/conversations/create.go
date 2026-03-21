package conversations

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
	Short: "Create a new conversation",
	RunE:  runCreate,
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

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	conversation, err := client.CreateConversation(context.Background(), clientAuth.ContactIdentifier)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(conversation)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
