package integrations

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

var hooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an integration hook",
	RunE:  runHooksCreate,
}

var hooksUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an integration hook",
	RunE:  runHooksUpdate,
}

var hooksDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an integration hook",
	RunE:  runHooksDelete,
}

func init() {
	hooksCreateCmd.Flags().String("app-id", "", "Integration app ID")
	hooksCreateCmd.MarkFlagRequired("app-id")
	hooksCreateCmd.Flags().Int("inbox-id", 0, "Inbox ID")
	hooksCreateCmd.Flags().String("settings", "", "Settings as JSON string")
	hooksCmd.AddCommand(hooksCreateCmd)

	hooksUpdateCmd.Flags().Int("hook-id", 0, "Hook ID")
	hooksUpdateCmd.MarkFlagRequired("hook-id")
	hooksUpdateCmd.Flags().String("settings", "", "Settings as JSON string")
	hooksCmd.AddCommand(hooksUpdateCmd)

	hooksDeleteCmd.Flags().Int("hook-id", 0, "Hook ID")
	hooksDeleteCmd.MarkFlagRequired("hook-id")
	hooksCmd.AddCommand(hooksDeleteCmd)
}

func runHooksCreate(cmd *cobra.Command, args []string) error {
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

	appID, _ := cmd.Flags().GetString("app-id")
	inboxID, _ := cmd.Flags().GetInt("inbox-id")
	opts := appapi.CreateIntegrationHookOpts{
		AppID:   appID,
		InboxID: inboxID,
	}

	if cmd.Flags().Changed("settings") {
		settingsStr, _ := cmd.Flags().GetString("settings")
		var settings any
		if err := json.Unmarshal([]byte(settingsStr), &settings); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeConfig, "invalid --settings JSON: "+err.Error())
		}
		opts.Settings = settings
	}

	ctx := context.Background()
	hook, err := client.CreateIntegrationHook(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(hook)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runHooksUpdate(cmd *cobra.Command, args []string) error {
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

	hookID, _ := cmd.Flags().GetInt("hook-id")
	opts := appapi.UpdateIntegrationHookOpts{}

	if cmd.Flags().Changed("settings") {
		settingsStr, _ := cmd.Flags().GetString("settings")
		var settings any
		if err := json.Unmarshal([]byte(settingsStr), &settings); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeConfig, "invalid --settings JSON: "+err.Error())
		}
		opts.Settings = settings
	}

	ctx := context.Background()
	hook, err := client.UpdateIntegrationHook(ctx, hookID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(hook)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runHooksDelete(cmd *cobra.Command, args []string) error {
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

	hookID, _ := cmd.Flags().GetInt("hook-id")
	ctx := context.Background()

	if err := client.DeleteIntegrationHook(ctx, hookID); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"deleted": true, "id": hookID})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
