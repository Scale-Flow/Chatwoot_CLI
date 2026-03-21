package application

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

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get authenticated user profile",
	RunE:  runProfileGet,
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update authenticated user profile",
	RunE:  runProfileUpdate,
}

func init() {
	profileUpdateCmd.Flags().String("name", "", "Display name")
	profileUpdateCmd.Flags().String("email", "", "Email address")
	profileUpdateCmd.Flags().String("availability", "", "Availability: online, offline, busy")
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileUpdateCmd)
}

func runProfileGet(cmd *cobra.Command, args []string) error {
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

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runProfileUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	availChanged := cmd.Flags().Changed("availability")
	if !nameChanged && !emailChanged && !availChanged {
		return fmt.Errorf("requires at least one of --name, --email, or --availability")
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

	opts := appapi.UpdateProfileOpts{}
	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if emailChanged {
		v, _ := cmd.Flags().GetString("email")
		opts.Email = &v
	}
	if availChanged {
		v, _ := cmd.Flags().GetString("availability")
		opts.Availability = &v
	}

	profile, err := client.UpdateProfile(context.Background(), opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
