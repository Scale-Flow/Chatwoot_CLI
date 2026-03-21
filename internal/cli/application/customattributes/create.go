package customattributes

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
	Short: "Create a custom attribute",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().String("attribute-key", "", "Attribute key")
	createCmd.MarkFlagRequired("attribute-key")
	createCmd.Flags().String("attribute-model", "", "Attribute model (contact or conversation)")
	createCmd.MarkFlagRequired("attribute-model")
	createCmd.Flags().String("attribute-type", "", "Attribute display type")
	createCmd.MarkFlagRequired("attribute-type")
	createCmd.Flags().String("name", "", "Attribute display name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("description", "", "Attribute description")
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

	key, _ := cmd.Flags().GetString("attribute-key")
	model, _ := cmd.Flags().GetString("attribute-model")
	attrType, _ := cmd.Flags().GetString("attribute-type")
	name, _ := cmd.Flags().GetString("name")
	desc, _ := cmd.Flags().GetString("description")

	opts := appapi.CreateCustomAttributeOpts{
		AttributeKey:         key,
		AttributeModel:       model,
		AttributeDisplayType: attrType,
		AttributeDisplayName: name,
		AttributeDescription: desc,
	}

	ctx := context.Background()
	attr, err := client.CreateCustomAttribute(ctx, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(attr)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
