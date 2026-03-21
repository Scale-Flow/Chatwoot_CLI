package conversations

import (
	"context"
	"encoding/json"
	"fmt"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter conversations by payload",
	RunE:  runFilter,
}

func init() {
	filterCmd.Flags().String("payload", "", "Filter payload (JSON array)")
	filterCmd.MarkFlagRequired("payload")
	filterCmd.Flags().Int("page", 1, "Page number")
	Cmd.AddCommand(filterCmd)
}

func runFilter(cmd *cobra.Command, args []string) error {
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

	payloadStr, _ := cmd.Flags().GetString("payload")
	var payload []any
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}

	page, _ := cmd.Flags().GetInt("page")
	ctx := context.Background()

	items, pag, err := client.FilterConversations(ctx, appapi.FilterConversationsOpts{
		Page:    page,
		Payload: payload,
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
