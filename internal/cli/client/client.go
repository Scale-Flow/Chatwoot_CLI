package client

import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/client/contacts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/client/conversations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/client/messages"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Client API commands (public end-user)",
}

func init() {
	Cmd.PersistentFlags().String("inbox-id", "", "Inbox identifier")
	Cmd.PersistentFlags().String("contact-id", "", "Contact identifier")
	Cmd.AddCommand(contacts.Cmd)
	Cmd.AddCommand(conversations.Cmd)
	Cmd.AddCommand(messages.Cmd)
}
