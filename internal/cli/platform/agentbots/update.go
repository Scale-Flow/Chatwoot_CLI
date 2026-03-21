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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a platform agent bot",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Agent bot ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Agent bot name")
	updateCmd.Flags().String("description", "", "Agent bot description")
	updateCmd.Flags().String("outgoing-url", "", "Outgoing URL for the agent bot")
	updateCmd.Flags().String("bot-type", "", "Bot type")
	updateCmd.Flags().String("bot-config", "", "Bot configuration as JSON")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	descChanged := cmd.Flags().Changed("description")
	urlChanged := cmd.Flags().Changed("outgoing-url")
	typeChanged := cmd.Flags().Changed("bot-type")
	cfgChanged := cmd.Flags().Changed("bot-config")

	if !nameChanged && !descChanged && !urlChanged && !typeChanged && !cfgChanged {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "at least one update flag required")
	}

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

	id, _ := cmd.Flags().GetInt("id")
	opts := platapi.UpdateAgentBotOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if descChanged {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = &v
	}
	if urlChanged {
		v, _ := cmd.Flags().GetString("outgoing-url")
		opts.OutgoingURL = &v
	}
	if typeChanged {
		v, _ := cmd.Flags().GetString("bot-type")
		opts.BotType = &v
	}
	if cfgChanged {
		v, _ := cmd.Flags().GetString("bot-config")
		var cfg any
		if err := json.Unmarshal([]byte(v), &cfg); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid JSON for --bot-config: "+err.Error())
		}
		opts.BotConfig = cfg
	}

	ctx := context.Background()
	bot, err := client.UpdateAgentBot(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
