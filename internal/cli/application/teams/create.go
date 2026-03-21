package teams

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
	Short: "Create a team",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Team name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("description", "", "Team description")
	createCmd.Flags().Bool("allow-auto-assign", false, "Allow auto assignment")
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
	opts := appapi.CreateTeamOpts{Name: name}

	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = v
	}
	if cmd.Flags().Changed("allow-auto-assign") {
		v, _ := cmd.Flags().GetBool("allow-auto-assign")
		opts.AllowAutoAssign = &v
	}

	ctx := context.Background()
	team, err := client.CreateTeam(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(team)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
