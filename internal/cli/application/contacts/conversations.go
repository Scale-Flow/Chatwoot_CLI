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

var conversationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contact conversations",
	RunE:  runConversationsList,
}

func init() {
	conversationsListCmd.Flags().Int("contact-id", 0, "Contact ID")
	conversationsListCmd.MarkFlagRequired("contact-id")
	contactConversationsCmd.AddCommand(conversationsListCmd)
}

func runConversationsList(cmd *cobra.Command, args []string) error {
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

	id, _ := cmd.Flags().GetInt("contact-id")
	ctx := context.Background()

	conversations, err := client.ListContactConversations(ctx, id)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(conversations, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
