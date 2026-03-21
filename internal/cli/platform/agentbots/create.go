package agentbots

import (
	"context"
	"encoding/json"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a platform agent bot",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Agent bot name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("description", "", "Agent bot description")
	createCmd.Flags().String("outgoing-url", "", "Outgoing URL for the agent bot")
	createCmd.Flags().String("bot-type", "", "Bot type")
	createCmd.Flags().String("bot-config", "", "Bot configuration as JSON")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}
	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModePlatform)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}
	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := platapi.NewClient(transport)

	name, _ := cmd.Flags().GetString("name")
	opts := platapi.CreateAgentBotOpts{Name: name}

	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = v
	}
	if cmd.Flags().Changed("outgoing-url") {
		v, _ := cmd.Flags().GetString("outgoing-url")
		opts.OutgoingURL = v
	}
	if cmd.Flags().Changed("bot-type") {
		v, _ := cmd.Flags().GetString("bot-type")
		opts.BotType = v
	}
	if cmd.Flags().Changed("bot-config") {
		v, _ := cmd.Flags().GetString("bot-config")
		var cfg any
		if err := json.Unmarshal([]byte(v), &cfg); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid JSON for --bot-config: "+err.Error())
		}
		opts.BotConfig = cfg
	}

	ctx := context.Background()
	bot, err := client.CreateAgentBot(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
