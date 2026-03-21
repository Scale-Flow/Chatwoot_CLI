package contacts

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	RunE:  runList,
}

func init() {
	cmdutil.AddPaginationFlags(listCmd)
	Cmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
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

	pf := cmdutil.GetPaginationFlags(cmd)
	ctx := context.Background()

	if pf.All {
		items, pag, err := chatwoot.ListAll(ctx, func(ctx context.Context, page int) ([]appapi.Contact, *contract.Pagination, error) {
			return client.ListContacts(ctx, appapi.ListContactsOpts{Page: page, PerPage: pf.PerPage})
		})
		if err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
		}
		meta := contract.Meta{}
		if pag != nil {
			meta.Pagination = pag
		}
		resp := contract.SuccessList(items, meta)
		return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
	}

	items, pag, err := client.ListContacts(ctx, appapi.ListContactsOpts{Page: pf.Page, PerPage: pf.PerPage})
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}
	meta := contract.Meta{}
	if pag != nil {
		meta.Pagination = pag
	}
	resp := contract.SuccessList(items, meta)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
