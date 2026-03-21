package webhooks

import (
	"context"
	"fmt"
	"strings"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a webhook",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Webhook ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("url", "", "Webhook URL")
	updateCmd.Flags().String("subscriptions", "", "Comma-separated event types")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	urlChanged := cmd.Flags().Changed("url")
	subsChanged := cmd.Flags().Changed("subscriptions")
	if !urlChanged && !subsChanged {
		return fmt.Errorf("requires at least one of --url or --subscriptions")
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
	opts := appapi.UpdateWebhookOpts{}

	if urlChanged {
		v, _ := cmd.Flags().GetString("url")
		opts.URL = &v
	}
	if subsChanged {
		subsStr, _ := cmd.Flags().GetString("subscriptions")
		if subsStr != "" {
			opts.Subscriptions = strings.Split(subsStr, ",")
		}
	}

	ctx := context.Background()
	webhook, err := client.UpdateWebhook(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(webhook)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
