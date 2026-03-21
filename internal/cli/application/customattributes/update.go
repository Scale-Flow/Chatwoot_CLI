package customattributes

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
	Short: "Update a custom attribute",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Custom attribute ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Attribute display name")
	updateCmd.Flags().String("description", "", "Attribute description")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	descChanged := cmd.Flags().Changed("description")
	if !nameChanged && !descChanged {
		return fmt.Errorf("requires at least one of --name or --description")
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
	opts := appapi.UpdateCustomAttributeOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.AttributeDisplayName = &v
	}
	if descChanged {
		v, _ := cmd.Flags().GetString("description")
		opts.AttributeDescription = &v
	}

	ctx := context.Background()
	attr, err := client.UpdateCustomAttribute(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(attr)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
