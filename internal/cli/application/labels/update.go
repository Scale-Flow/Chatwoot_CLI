package labels

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
	Short: "Update a label",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Label ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("title", "", "Label title")
	updateCmd.Flags().String("description", "", "Label description")
	updateCmd.Flags().String("color", "", "Label color")
	updateCmd.Flags().Bool("show-on-sidebar", false, "Show label on sidebar")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	titleChanged := cmd.Flags().Changed("title")
	descChanged := cmd.Flags().Changed("description")
	colorChanged := cmd.Flags().Changed("color")
	sidebarChanged := cmd.Flags().Changed("show-on-sidebar")
	if !titleChanged && !descChanged && !colorChanged && !sidebarChanged {
		return fmt.Errorf("requires at least one of --title, --description, --color, or --show-on-sidebar")
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
	opts := appapi.UpdateLabelOpts{}

	if titleChanged {
		v, _ := cmd.Flags().GetString("title")
		opts.Title = &v
	}
	if descChanged {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = &v
	}
	if colorChanged {
		v, _ := cmd.Flags().GetString("color")
		opts.Color = &v
	}
	if sidebarChanged {
		v, _ := cmd.Flags().GetBool("show-on-sidebar")
		opts.ShowOnSidebar = &v
	}

	ctx := context.Background()
	label, err := client.UpdateLabel(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(label)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
