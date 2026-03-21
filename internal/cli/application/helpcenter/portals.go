package helpcenter

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

var portalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List portals",
	RunE:  runPortalsList,
}

var portalsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a portal",
	RunE:  runPortalsCreate,
}

var portalsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a portal",
	RunE:  runPortalsUpdate,
}

func init() {
	portalsCmd.AddCommand(portalsListCmd)

	portalsCreateCmd.Flags().String("name", "", "Portal name")
	portalsCreateCmd.MarkFlagRequired("name")
	portalsCreateCmd.Flags().String("slug", "", "Portal slug")
	portalsCreateCmd.Flags().String("color", "", "Portal color")
	portalsCreateCmd.Flags().String("header-text", "", "Portal header text")
	portalsCreateCmd.Flags().String("custom-domain", "", "Portal custom domain")
	portalsCmd.AddCommand(portalsCreateCmd)

	portalsUpdateCmd.Flags().Int("id", 0, "Portal ID")
	portalsUpdateCmd.MarkFlagRequired("id")
	portalsUpdateCmd.Flags().String("name", "", "Portal name")
	portalsUpdateCmd.Flags().String("slug", "", "Portal slug")
	portalsUpdateCmd.Flags().String("color", "", "Portal color")
	portalsUpdateCmd.Flags().String("header-text", "", "Portal header text")
	portalsUpdateCmd.Flags().String("custom-domain", "", "Portal custom domain")
	portalsUpdateCmd.Flags().Bool("archived", false, "Archive portal")
	portalsCmd.AddCommand(portalsUpdateCmd)
}

func runPortalsList(cmd *cobra.Command, args []string) error {
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

	ctx := context.Background()
	portals, err := client.ListPortals(ctx)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(portals, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runPortalsCreate(cmd *cobra.Command, args []string) error {
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
	opts := appapi.CreatePortalOpts{Name: name}

	if cmd.Flags().Changed("slug") {
		v, _ := cmd.Flags().GetString("slug")
		opts.Slug = v
	}
	if cmd.Flags().Changed("color") {
		v, _ := cmd.Flags().GetString("color")
		opts.Color = v
	}
	if cmd.Flags().Changed("header-text") {
		v, _ := cmd.Flags().GetString("header-text")
		opts.HeaderText = v
	}
	if cmd.Flags().Changed("custom-domain") {
		v, _ := cmd.Flags().GetString("custom-domain")
		opts.CustomDomain = v
	}

	ctx := context.Background()
	portal, err := client.CreatePortal(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(portal)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runPortalsUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	slugChanged := cmd.Flags().Changed("slug")
	colorChanged := cmd.Flags().Changed("color")
	headerChanged := cmd.Flags().Changed("header-text")
	domainChanged := cmd.Flags().Changed("custom-domain")
	archivedChanged := cmd.Flags().Changed("archived")

	if !nameChanged && !slugChanged && !colorChanged && !headerChanged && !domainChanged && !archivedChanged {
		return fmt.Errorf("requires at least one of --name, --slug, --color, --header-text, --custom-domain, or --archived")
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
	opts := appapi.UpdatePortalOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if slugChanged {
		v, _ := cmd.Flags().GetString("slug")
		opts.Slug = &v
	}
	if colorChanged {
		v, _ := cmd.Flags().GetString("color")
		opts.Color = &v
	}
	if headerChanged {
		v, _ := cmd.Flags().GetString("header-text")
		opts.HeaderText = &v
	}
	if domainChanged {
		v, _ := cmd.Flags().GetString("custom-domain")
		opts.CustomDomain = &v
	}
	if archivedChanged {
		v, _ := cmd.Flags().GetBool("archived")
		opts.Archived = &v
	}

	ctx := context.Background()
	portal, err := client.UpdatePortal(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(portal)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
