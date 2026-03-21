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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a custom filter",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("name", "", "Filter name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("type", "", "Filter type (conversation, contact, report)")
	createCmd.MarkFlagRequired("type")
	createCmd.Flags().String("query", "", "Filter query as JSON string")
	createCmd.MarkFlagRequired("query")
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

	name, _ := cmd.Flags().GetString("name")
	filterType, _ := cmd.Flags().GetString("type")
	queryStr, _ := cmd.Flags().GetString("query")

	var query any
	if err := json.Unmarshal([]byte(queryStr), &query); err != nil {
		return fmt.Errorf("invalid --query JSON: %w", err)
	}

	opts := appapi.CreateCustomFilterOpts{
		Name:  name,
		Type:  filterType,
		Query: query,
	}

	ctx := context.Background()
	filter, err := client.CreateCustomFilter(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(filter)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
