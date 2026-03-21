package accountusers

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a user from an account",
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().Int("user-id", 0, "User ID")
	deleteCmd.MarkFlagRequired("user-id")
	Cmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
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

	userID, _ := cmd.Flags().GetInt("user-id")

	ctx := context.Background()
	if err := client.DeleteAccountUser(ctx, rctx.AccountID, userID); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"deleted": true})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
