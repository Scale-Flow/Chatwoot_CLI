package labels

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
	Short: "Create a label",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("title", "", "Label title")
	createCmd.MarkFlagRequired("title")
	createCmd.Flags().String("description", "", "Label description")
	createCmd.Flags().String("color", "", "Label color")
	createCmd.Flags().Bool("show-on-sidebar", false, "Show label on sidebar")
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

	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	color, _ := cmd.Flags().GetString("color")

	opts := appapi.CreateLabelOpts{
		Title:       title,
		Description: description,
		Color:       color,
	}

	if cmd.Flags().Changed("show-on-sidebar") {
		v, _ := cmd.Flags().GetBool("show-on-sidebar")
		opts.ShowOnSidebar = &v
	}

	ctx := context.Background()
	label, err := client.CreateLabel(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(label)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
