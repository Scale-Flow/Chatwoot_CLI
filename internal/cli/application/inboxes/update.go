package inboxes

import (
	"context"
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
	Short: "Update an inbox",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Inbox ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Inbox name")
	updateCmd.Flags().Bool("enable-auto-assignment", false, "Enable auto assignment")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	autoChanged := cmd.Flags().Changed("enable-auto-assignment")
	if !nameChanged && !autoChanged {
		return fmt.Errorf("requires at least one of --name or --enable-auto-assignment")
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
	opts := appapi.UpdateInboxOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if autoChanged {
		v, _ := cmd.Flags().GetBool("enable-auto-assignment")
		opts.EnableAutoAssignment = &v
	}

	ctx := context.Background()
	inbox, err := client.UpdateInbox(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(inbox)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
