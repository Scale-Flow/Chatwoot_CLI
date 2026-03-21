package agentbots

import (
	"context"
	"encoding/json"
	"fmt"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an agent bot",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Agent bot ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Agent bot name")
	updateCmd.Flags().String("description", "", "Agent bot description")
	updateCmd.Flags().String("bot-type", "", "Agent bot type")
	updateCmd.Flags().String("outgoing-url", "", "Agent bot outgoing URL")
	updateCmd.Flags().String("bot-config", "", "Agent bot config as JSON string")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	descChanged := cmd.Flags().Changed("description")
	botTypeChanged := cmd.Flags().Changed("bot-type")
	urlChanged := cmd.Flags().Changed("outgoing-url")
	configChanged := cmd.Flags().Changed("bot-config")
	if !nameChanged && !descChanged && !botTypeChanged && !urlChanged && !configChanged {
		return fmt.Errorf("requires at least one of --name, --description, --bot-type, --outgoing-url, or --bot-config")
	}

	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	id, _ := cmd.Flags().GetInt("id")
	opts := appapi.UpdateAgentBotOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if descChanged {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = &v
	}
	if botTypeChanged {
		v, _ := cmd.Flags().GetString("bot-type")
		opts.BotType = &v
	}
	if urlChanged {
		v, _ := cmd.Flags().GetString("outgoing-url")
		opts.OutgoingURL = &v
	}
	if configChanged {
		configStr, _ := cmd.Flags().GetString("bot-config")
		var config any
		if err := json.Unmarshal([]byte(configStr), &config); err != nil {
			return fmt.Errorf("invalid --bot-config JSON: %w", err)
		}
		opts.BotConfig = config
	}

	ctx := context.Background()
	bot, err := client.UpdateAgentBot(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
