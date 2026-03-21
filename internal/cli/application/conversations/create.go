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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a conversation",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().Int("contact-id", 0, "Contact ID")
	createCmd.MarkFlagRequired("contact-id")
	createCmd.Flags().Int("inbox-id", 0, "Inbox ID")
	createCmd.MarkFlagRequired("inbox-id")
	createCmd.Flags().String("status", "", "Conversation status")
	createCmd.Flags().String("content", "", "Initial message content")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
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

	contactID, _ := cmd.Flags().GetInt("contact-id")
	inboxID, _ := cmd.Flags().GetInt("inbox-id")
	status, _ := cmd.Flags().GetString("status")
	content, _ := cmd.Flags().GetString("content")

	opts := appapi.CreateConversationOpts{
		ContactID: contactID,
		InboxID:   inboxID,
	}
	if status != "" {
		opts.Status = status
	}
	if content != "" {
		opts.Message = &struct {
			Content string `json:"content"`
		}{Content: content}
	}

	ctx := context.Background()
	convo, err := client.CreateConversation(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(convo)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
