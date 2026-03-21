package conversations

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List conversations for a contact",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
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

	items, err := client.ListConversations(context.Background(), clientAuth.ContactIdentifier)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(items, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
