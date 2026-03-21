package contacts

import (
	"context"
	"strings"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var labelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contact labels",
	RunE:  runLabelsList,
}

var labelsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set contact labels",
	RunE:  runLabelsSet,
}

func init() {
	labelsListCmd.Flags().Int("id", 0, "Contact ID")
	labelsListCmd.MarkFlagRequired("id")
	labelsCmd.AddCommand(labelsListCmd)

	labelsSetCmd.Flags().Int("id", 0, "Contact ID")
	labelsSetCmd.MarkFlagRequired("id")
	labelsSetCmd.Flags().String("labels", "", "Comma-separated labels")
	labelsSetCmd.MarkFlagRequired("labels")
	labelsCmd.AddCommand(labelsSetCmd)
}

func runLabelsList(cmd *cobra.Command, args []string) error {
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
	ctx := context.Background()

	labels, err := client.ListContactLabels(ctx, id)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(labels)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runLabelsSet(cmd *cobra.Command, args []string) error {
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
	labelsStr, _ := cmd.Flags().GetString("labels")
	labels := strings.Split(labelsStr, ",")
	for i := range labels {
		labels[i] = strings.TrimSpace(labels[i])
	}

	ctx := context.Background()

	result, err := client.SetContactLabels(ctx, id, labels)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
