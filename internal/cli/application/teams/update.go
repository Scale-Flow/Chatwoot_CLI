package teams

import (
	"context"
	"fmt"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a team",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Team ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Team name")
	updateCmd.Flags().String("description", "", "Team description")
	updateCmd.Flags().Bool("allow-auto-assign", false, "Allow auto assignment")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("description") && !cmd.Flags().Changed("allow-auto-assign") {
		return fmt.Errorf("requires at least one of --name, --description, or --allow-auto-assign")
	}

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
	opts := appapi.UpdateTeamOpts{}

	if cmd.Flags().Changed("name") {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = &v
	}
	if cmd.Flags().Changed("allow-auto-assign") {
		v, _ := cmd.Flags().GetBool("allow-auto-assign")
		opts.AllowAutoAssign = &v
	}

	ctx := context.Background()
	team, err := client.UpdateTeam(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(team)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
