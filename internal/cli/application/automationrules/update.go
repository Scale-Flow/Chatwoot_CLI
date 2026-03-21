package automationrules

import (
	"context"
	"encoding/json"
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
	Short: "Update an automation rule",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Automation rule ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Rule name")
	updateCmd.Flags().String("event-name", "", "Event name that triggers the rule")
	updateCmd.Flags().String("conditions", "", "Conditions as JSON array")
	updateCmd.Flags().String("actions", "", "Actions as JSON array")
	updateCmd.Flags().String("description", "", "Rule description")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	eventNameChanged := cmd.Flags().Changed("event-name")
	conditionsChanged := cmd.Flags().Changed("conditions")
	actionsChanged := cmd.Flags().Changed("actions")
	descriptionChanged := cmd.Flags().Changed("description")

	if !nameChanged && !eventNameChanged && !conditionsChanged && !actionsChanged && !descriptionChanged {
		return fmt.Errorf("requires at least one of --name, --event-name, --conditions, --actions, or --description")
	}

	opts := appapi.UpdateAutomationRuleOpts{}

	if conditionsChanged {
		conditionsStr, _ := cmd.Flags().GetString("conditions")
		var conditions any
		if err := json.Unmarshal([]byte(conditionsStr), &conditions); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid conditions JSON: "+err.Error())
		}
		opts.Conditions = conditions
	}

	if actionsChanged {
		actionsStr, _ := cmd.Flags().GetString("actions")
		var actions any
		if err := json.Unmarshal([]byte(actionsStr), &actions); err != nil {
			return cmdutil.WriteError(cmd, contract.ErrCodeValidation, "invalid actions JSON: "+err.Error())
		}
		opts.Actions = actions
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

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if eventNameChanged {
		v, _ := cmd.Flags().GetString("event-name")
		opts.EventName = &v
	}
	if descriptionChanged {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = &v
	}

	ctx := context.Background()
	rule, err := client.UpdateAutomationRule(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(rule)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
