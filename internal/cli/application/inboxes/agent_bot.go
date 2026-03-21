package inboxes

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var agentBotGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get inbox agent bot",
	RunE:  runAgentBotGet,
}

var agentBotSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set inbox agent bot",
	RunE:  runAgentBotSet,
}

func init() {
	agentBotGetCmd.Flags().Int("inbox-id", 0, "Inbox ID")
	agentBotGetCmd.MarkFlagRequired("inbox-id")
	agentBotCmd.AddCommand(agentBotGetCmd)

	agentBotSetCmd.Flags().Int("inbox-id", 0, "Inbox ID")
	agentBotSetCmd.MarkFlagRequired("inbox-id")
	agentBotSetCmd.Flags().Int("agent-bot-id", 0, "Agent bot ID")
	agentBotSetCmd.MarkFlagRequired("agent-bot-id")
	agentBotCmd.AddCommand(agentBotSetCmd)
}

func runAgentBotGet(cmd *cobra.Command, args []string) error {
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

	inboxID, _ := cmd.Flags().GetInt("inbox-id")
	ctx := context.Background()

	bot, err := client.GetInboxAgentBot(ctx, inboxID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runAgentBotSet(cmd *cobra.Command, args []string) error {
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

	inboxID, _ := cmd.Flags().GetInt("inbox-id")
	botID, _ := cmd.Flags().GetInt("agent-bot-id")
	ctx := context.Background()

	bot, err := client.SetInboxAgentBot(ctx, inboxID, botID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(bot)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
