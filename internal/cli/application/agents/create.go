package agents

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
	Short: "Create an agent",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Agent name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("email", "", "Agent email")
	createCmd.MarkFlagRequired("email")
	createCmd.Flags().String("role", "agent", "Agent role")
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
	email, _ := cmd.Flags().GetString("email")
	role, _ := cmd.Flags().GetString("role")

	opts := appapi.CreateAgentOpts{
		Name:  name,
		Email: email,
		Role:  role,
	}

	ctx := context.Background()
	agent, err := client.CreateAgent(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(agent)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
