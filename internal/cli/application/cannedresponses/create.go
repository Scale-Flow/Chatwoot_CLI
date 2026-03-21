package cannedresponses

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a canned response",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("short-code", "", "Short code for the canned response")
	createCmd.MarkFlagRequired("short-code")
	createCmd.Flags().String("content", "", "Content of the canned response")
	createCmd.MarkFlagRequired("content")
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

	shortCode, _ := cmd.Flags().GetString("short-code")
	content, _ := cmd.Flags().GetString("content")

	opts := appapi.CreateCannedResponseOpts{
		ShortCode: shortCode,
		Content:   content,
	}

	ctx := context.Background()
	result, err := client.CreateCannedResponse(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
