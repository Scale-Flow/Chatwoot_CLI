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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an agent bot",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Agent bot name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("description", "", "Agent bot description")
	createCmd.Flags().String("bot-type", "", "Agent bot type")
	createCmd.Flags().String("outgoing-url", "", "Agent bot outgoing URL")
	createCmd.Flags().String("bot-config", "", "Agent bot config as JSON string")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
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

	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	botType, _ := cmd.Flags().GetString("bot-type")
	outgoingURL, _ := cmd.Flags().GetString("outgoing-url")

	opts := appapi.CreateAgentBotOpts{
		Name:        name,
		Description: description,
		BotType:     botType,
		OutgoingURL: outgoingURL,
	}

	if cmd.Flags().Changed("bot-config") {
		configStr, _ := cmd.Flags().GetString("bot-config")
		var config any
		if err := json.Unmarshal([]byte(configStr), &config); err != nil {
			return fmt.Errorf("invalid --bot-config JSON: %w", err)
		}
		opts.BotConfig = config
	}

	ctx := context.Background()
	bot, err := client.CreateAgentBot(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
