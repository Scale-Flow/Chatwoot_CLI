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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a platform user",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "User ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "User name")
	updateCmd.Flags().String("email", "", "User email")
	updateCmd.Flags().String("password", "", "User password")
	updateCmd.Flags().String("custom-attributes", "", "Custom attributes as JSON")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	passwordChanged := cmd.Flags().Changed("password")
	attrsChanged := cmd.Flags().Changed("custom-attributes")

	if !nameChanged && !emailChanged && !passwordChanged && !attrsChanged {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "at least one update flag required")
	}

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

	id, _ := cmd.Flags().GetInt("id")
	opts := platapi.UpdateUserOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if emailChanged {
		v, _ := cmd.Flags().GetString("email")
		opts.Email = &v
	}
	if passwordChanged {
		v, _ := cmd.Flags().GetString("password")
		opts.Password = &v
	}
	if attrsChanged {
		v, _ := cmd.Flags().GetString("custom-attributes")
		var attrs any
		if err := json.Unmarshal([]byte(v), &attrs); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid JSON for --custom-attributes: "+err.Error())
		}
		opts.CustomAttributes = attrs
	}

	ctx := context.Background()
	user, err := client.UpdateUser(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(user)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
