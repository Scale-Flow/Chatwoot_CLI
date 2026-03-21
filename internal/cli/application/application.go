package application

import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/account"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/agentbots"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/agents"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/auditlogs"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/automationrules"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/cannedresponses"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/contacts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/conversations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/customattributes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/customfilters"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/helpcenter"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/inboxes"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/integrations"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/labels"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/messages"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/reports"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/teams"
	"github.com/chatwoot/chatwoot-cli/internal/cli/application/webhooks"
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
	// Sprint C
	Cmd.AddCommand(contacts.Cmd)
	Cmd.AddCommand(conversations.Cmd)
	Cmd.AddCommand(messages.Cmd)
	Cmd.AddCommand(inboxes.Cmd)
	// Sprint D
	Cmd.AddCommand(teams.Cmd)
	Cmd.AddCommand(agents.Cmd)
	Cmd.AddCommand(cannedresponses.Cmd)
	Cmd.AddCommand(reports.Cmd)
	Cmd.AddCommand(webhooks.Cmd)
	Cmd.AddCommand(automationrules.Cmd)
	Cmd.AddCommand(labels.Cmd)
	Cmd.AddCommand(customattributes.Cmd)
	Cmd.AddCommand(customfilters.Cmd)
	Cmd.AddCommand(account.Cmd)
	Cmd.AddCommand(agentbots.Cmd)
	Cmd.AddCommand(auditlogs.Cmd)
	Cmd.AddCommand(integrations.Cmd)
	Cmd.AddCommand(helpcenter.Cmd)
}
