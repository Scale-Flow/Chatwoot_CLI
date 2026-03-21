package conversations

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
	Short: "List conversations",
	RunE:  runList,
}

func init() {
	cmdutil.AddPaginationFlags(listCmd)
	listCmd.Flags().String("status", "", "Filter by status (open, resolved, pending, snoozed)")
	listCmd.Flags().Int("inbox-id", 0, "Filter by inbox ID")
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
	status, _ := cmd.Flags().GetString("status")
	inboxID, _ := cmd.Flags().GetInt("inbox-id")
	ctx := context.Background()

	if pf.All {
		items, pag, err := chatwoot.ListAll(ctx, func(ctx context.Context, page int) ([]appapi.Conversation, *contract.Pagination, error) {
			return client.ListConversations(ctx, appapi.ListConversationsOpts{
				Page:    page,
				PerPage: pf.PerPage,
				Status:  status,
				InboxID: inboxID,
			})
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

	items, pag, err := client.ListConversations(ctx, appapi.ListConversationsOpts{
		Page:    pf.Page,
		PerPage: pf.PerPage,
		Status:  status,
		InboxID: inboxID,
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
