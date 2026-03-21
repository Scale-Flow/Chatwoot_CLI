package helpcenter

import (
	"context"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var categoriesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a category",
	RunE:  runCategoriesCreate,
}

func init() {
	categoriesCreateCmd.Flags().Int("portal-id", 0, "Portal ID")
	categoriesCreateCmd.MarkFlagRequired("portal-id")
	categoriesCreateCmd.Flags().String("name", "", "Category name")
	categoriesCreateCmd.MarkFlagRequired("name")
	categoriesCreateCmd.Flags().String("description", "", "Category description")
	categoriesCreateCmd.Flags().String("locale", "", "Category locale")
	categoriesCreateCmd.Flags().Int("position", 0, "Category position")
	categoriesCmd.AddCommand(categoriesCreateCmd)
}

func runCategoriesCreate(cmd *cobra.Command, args []string) error {
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

	portalID, _ := cmd.Flags().GetInt("portal-id")
	name, _ := cmd.Flags().GetString("name")
	opts := appapi.CreateCategoryOpts{Name: name}

	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = v
	}
	if cmd.Flags().Changed("locale") {
		v, _ := cmd.Flags().GetString("locale")
		opts.Locale = v
	}
	if cmd.Flags().Changed("position") {
		v, _ := cmd.Flags().GetInt("position")
		opts.Position = v
	}

	ctx := context.Background()
	category, err := client.CreateCategory(ctx, portalID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(category)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
