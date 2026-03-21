package webhooks

import (
	"context"
	"strings"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("url", "", "Webhook URL")
	createCmd.MarkFlagRequired("url")
	createCmd.Flags().String("subscriptions", "", "Comma-separated event types")
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

	url, _ := cmd.Flags().GetString("url")
	subsStr, _ := cmd.Flags().GetString("subscriptions")
	var subs []string
	if subsStr != "" {
		subs = strings.Split(subsStr, ",")
	}
	opts := appapi.CreateWebhookOpts{URL: url, Subscriptions: subs}

	ctx := context.Background()
	webhook, err := client.CreateWebhook(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(webhook)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
