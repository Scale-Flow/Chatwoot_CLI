package contacts

import (
	"context"

	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	clientapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/clientapi"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a client contact",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().String("name", "", "Contact name")
	updateCmd.Flags().String("email", "", "Contact email")
	updateCmd.Flags().String("phone", "", "Contact phone number")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}
	clientAuth, err := cmdutil.ResolveClientAuth(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}
	if clientAuth.ContactIdentifier == "" {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, "--contact-id flag or CHATWOOT_CONTACT_IDENTIFIER env var required for this command")
	}

	// Require at least one update flag
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	phoneChanged := cmd.Flags().Changed("phone")
	if !nameChanged && !emailChanged && !phoneChanged {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "at least one update flag required")
	}

	opts := clientapi.UpdateContactOpts{}
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

	transport := chatwoot.NewClient(rctx.BaseURL, "", "")
	client := clientapi.NewClient(transport, clientAuth.InboxIdentifier)

	contact, err := client.UpdateContact(context.Background(), clientAuth.ContactIdentifier, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(contact)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
