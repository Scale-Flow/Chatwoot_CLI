package account

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
	Short: "Update account settings",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().String("name", "", "Account name")
	updateCmd.Flags().String("locale", "", "Account locale")
	updateCmd.Flags().String("domain", "", "Account domain")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	localeChanged := cmd.Flags().Changed("locale")
	domainChanged := cmd.Flags().Changed("domain")
	if !nameChanged && !localeChanged && !domainChanged {
		return fmt.Errorf("requires at least one of --name, --locale, or --domain")
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

	opts := appapi.UpdateAccountOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if localeChanged {
		v, _ := cmd.Flags().GetString("locale")
		opts.Locale = &v
	}
	if domainChanged {
		v, _ := cmd.Flags().GetString("domain")
		opts.Domain = &v
	}

	ctx := context.Background()
	info, err := client.UpdateAccount(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(info)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
