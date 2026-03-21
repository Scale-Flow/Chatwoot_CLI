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

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search contacts by query",
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().String("query", "", "Search query")
	searchCmd.MarkFlagRequired("query")
	searchCmd.Flags().Int("page", 1, "Page number")
	searchCmd.Flags().Int("per-page", 25, "Items per page")
	Cmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
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

	query, _ := cmd.Flags().GetString("query")
	page, _ := cmd.Flags().GetInt("page")
	ctx := context.Background()

	items, pag, err := client.SearchContacts(ctx, query, page)
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
