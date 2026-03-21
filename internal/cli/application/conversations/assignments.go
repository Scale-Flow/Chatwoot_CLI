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

var assignmentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Assign a conversation",
	RunE:  runAssignmentsCreate,
}

func init() {
	assignmentsCreateCmd.Flags().Int("id", 0, "Conversation ID")
	assignmentsCreateCmd.MarkFlagRequired("id")
	assignmentsCreateCmd.Flags().Int("agent-id", 0, "Agent ID to assign")
	assignmentsCreateCmd.Flags().Int("team-id", 0, "Team ID to assign")
	assignmentsCmd.AddCommand(assignmentsCreateCmd)
}

func runAssignmentsCreate(cmd *cobra.Command, args []string) error {
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
	opts := appapi.AssignOpts{}

	if cmd.Flags().Changed("agent-id") {
		v, _ := cmd.Flags().GetInt("agent-id")
		opts.AgentID = &v
	}
	if cmd.Flags().Changed("team-id") {
		v, _ := cmd.Flags().GetInt("team-id")
		opts.TeamID = &v
	}

	ctx := context.Background()
	assignment, err := client.AssignConversation(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(assignment)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
