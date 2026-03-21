package inboxes

import (
	"context"
	"encoding/json"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an inbox",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Inbox name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("channel", "", "Channel configuration (JSON string)")
	createCmd.MarkFlagRequired("channel")
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

	name, _ := cmd.Flags().GetString("name")
	channelStr, _ := cmd.Flags().GetString("channel")

	var channel any
	if err := json.Unmarshal([]byte(channelStr), &channel); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid --channel JSON: "+err.Error())
	}

	ctx := context.Background()
	inbox, err := client.CreateInbox(ctx, appapi.CreateInboxOpts{Name: name, Channel: channel})
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(inbox)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
