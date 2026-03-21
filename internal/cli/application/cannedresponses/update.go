package cannedresponses

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
	Short: "Update a canned response",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Canned response ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("short-code", "", "Short code for the canned response")
	updateCmd.Flags().String("content", "", "Content of the canned response")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	shortCodeChanged := cmd.Flags().Changed("short-code")
	contentChanged := cmd.Flags().Changed("content")
	if !shortCodeChanged && !contentChanged {
		return fmt.Errorf("requires at least one of --short-code or --content")
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
	opts := appapi.UpdateCannedResponseOpts{}

	if shortCodeChanged {
		v, _ := cmd.Flags().GetString("short-code")
		opts.ShortCode = &v
	}
	if contentChanged {
		v, _ := cmd.Flags().GetString("content")
		opts.Content = &v
	}

	ctx := context.Background()
	result, err := client.UpdateCannedResponse(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
