package users

import (
	"context"
	"encoding/json"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	platapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/platform"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a platform user",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "User name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("email", "", "User email")
	createCmd.MarkFlagRequired("email")
	createCmd.Flags().String("password", "", "User password")
	createCmd.MarkFlagRequired("password")
	createCmd.Flags().String("custom-attributes", "", "Custom attributes as JSON")
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

	name, _ := cmd.Flags().GetString("name")
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")
	opts := platapi.CreateUserOpts{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if cmd.Flags().Changed("custom-attributes") {
		v, _ := cmd.Flags().GetString("custom-attributes")
		var attrs any
		if err := json.Unmarshal([]byte(v), &attrs); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid JSON for --custom-attributes: "+err.Error())
		}
		opts.CustomAttributes = attrs
	}

	ctx := context.Background()
	user, err := client.CreateUser(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(user)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
