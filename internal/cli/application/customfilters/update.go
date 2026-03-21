package customfilters

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
	Short: "Update a custom filter",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().Int("id", 0, "Custom filter ID")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("name", "", "Filter name")
	updateCmd.Flags().String("query", "", "Filter query as JSON string")
	Cmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	nameChanged := cmd.Flags().Changed("name")
	queryChanged := cmd.Flags().Changed("query")
	if !nameChanged && !queryChanged {
		return fmt.Errorf("requires at least one of --name or --query")
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
	opts := appapi.UpdateCustomFilterOpts{}

	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if queryChanged {
		queryStr, _ := cmd.Flags().GetString("query")
		var query any
		if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
			return fmt.Errorf("invalid --query JSON: %w", err)
		}
		opts.Query = query
	}

	ctx := context.Background()
	filter, err := client.UpdateCustomFilter(ctx, id, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(filter)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
