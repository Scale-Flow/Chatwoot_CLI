package automationrules

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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an automation rule",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Rule name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("event-name", "", "Event name that triggers the rule")
	createCmd.MarkFlagRequired("event-name")
	createCmd.Flags().String("conditions", "", "Conditions as JSON array")
	createCmd.MarkFlagRequired("conditions")
	createCmd.Flags().String("actions", "", "Actions as JSON array")
	createCmd.MarkFlagRequired("actions")
	createCmd.Flags().String("description", "", "Rule description")
	Cmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	conditionsStr, _ := cmd.Flags().GetString("conditions")
	var conditions any
	if err := json.Unmarshal([]byte(conditionsStr), &conditions); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid conditions JSON: "+err.Error())
	}

	actionsStr, _ := cmd.Flags().GetString("actions")
	var actions any
	if err := json.Unmarshal([]byte(actionsStr), &actions); err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid actions JSON: "+err.Error())
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

	name, _ := cmd.Flags().GetString("name")
	eventName, _ := cmd.Flags().GetString("event-name")
	description, _ := cmd.Flags().GetString("description")

	opts := appapi.CreateAutomationRuleOpts{
		Name:        name,
		EventName:   eventName,
		Conditions:  conditions,
		Actions:     actions,
		Description: description,
	}

	ctx := context.Background()
	rule, err := client.CreateAutomationRule(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(rule)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
