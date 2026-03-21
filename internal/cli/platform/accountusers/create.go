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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a user to an account",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().Int("user-id", 0, "User ID")
	createCmd.MarkFlagRequired("user-id")
	createCmd.Flags().String("role", "", "Role for the user in the account")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
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
	role, _ := cmd.Flags().GetString("role")
	opts := platapi.CreateAccountUserOpts{
		UserID: userID,
		Role:   role,
	}

	ctx := context.Background()
	if err := client.CreateAccountUser(ctx, rctx.AccountID, opts); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"created": true})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
