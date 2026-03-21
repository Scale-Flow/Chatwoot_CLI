package platform

import (
	"github.com/chatwoot/chatwoot-cli/internal/cli/platform/accounts"
	"github.com/chatwoot/chatwoot-cli/internal/cli/platform/accountusers"
	"github.com/chatwoot/chatwoot-cli/internal/cli/platform/agentbots"
	"github.com/chatwoot/chatwoot-cli/internal/cli/platform/users"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "platform",
	Short: "Platform API commands (self-hosted admin)",
}

func init() {
	Cmd.AddCommand(accounts.Cmd)
	Cmd.AddCommand(accountusers.Cmd)
	Cmd.AddCommand(agentbots.Cmd)
	Cmd.AddCommand(users.Cmd)
}
