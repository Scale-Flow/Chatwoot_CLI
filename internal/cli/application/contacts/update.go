package contacts

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
	Short: "Update a contact",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Contact ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Contact name")
	updateCmd.Flags().String("email", "", "Contact email")
	updateCmd.Flags().String("phone", "", "Contact phone number")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	phoneChanged := cmd.Flags().Changed("phone")
	if !nameChanged && !emailChanged && !phoneChanged {
		return fmt.Errorf("requires at least one of --name, --email, or --phone")
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
	opts := appapi.UpdateContactOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if emailChanged {
		v, _ := cmd.Flags().GetString("email")
		opts.Email = &v
	}
	if phoneChanged {
		v, _ := cmd.Flags().GetString("phone")
		opts.Phone = &v
	}

	ctx := context.Background()
	contact, err := client.UpdateContact(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(contact)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
