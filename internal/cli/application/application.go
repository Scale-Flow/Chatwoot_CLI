package application

import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/contacts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/conversations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/inboxes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/messages"
	"github.com/spf13/cobra"
)

// Cmd is the application command group.
var Cmd = &cobra.Command{
	Use:   "application",
	Short: "Application API commands (agent/admin)",
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage authenticated user profile",
}

func init() {
	Cmd.AddCommand(profileCmd)
	Cmd.AddCommand(contacts.Cmd)
	Cmd.AddCommand(conversations.Cmd)
	Cmd.AddCommand(messages.Cmd)
	Cmd.AddCommand(inboxes.Cmd)
}
