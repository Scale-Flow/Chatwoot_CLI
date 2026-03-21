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

var articlesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an article",
	RunE:  runArticlesCreate,
}

func init() {
	articlesCreateCmd.Flags().Int("portal-id", 0, "Portal ID")
	articlesCreateCmd.MarkFlagRequired("portal-id")
	articlesCreateCmd.Flags().String("title", "", "Article title")
	articlesCreateCmd.MarkFlagRequired("title")
	articlesCreateCmd.Flags().String("content", "", "Article content")
	articlesCreateCmd.Flags().String("description", "", "Article description")
	articlesCreateCmd.Flags().Int("status", 0, "Article status")
	articlesCreateCmd.Flags().Int("category-id", 0, "Category ID")
	articlesCreateCmd.Flags().Int("author-id", 0, "Author ID")
	articlesCmd.AddCommand(articlesCreateCmd)
}

func runArticlesCreate(cmd *cobra.Command, args []string) error {
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
	title, _ := cmd.Flags().GetString("title")
	opts := appapi.CreateArticleOpts{Title: title}

	if cmd.Flags().Changed("content") {
		v, _ := cmd.Flags().GetString("content")
		opts.Content = v
	}
	if cmd.Flags().Changed("description") {
		v, _ := cmd.Flags().GetString("description")
		opts.Description = v
	}
	if cmd.Flags().Changed("status") {
		v, _ := cmd.Flags().GetInt("status")
		opts.Status = v
	}
	if cmd.Flags().Changed("category-id") {
		v, _ := cmd.Flags().GetInt("category-id")
		opts.CategoryID = v
	}
	if cmd.Flags().Changed("author-id") {
		v, _ := cmd.Flags().GetInt("author-id")
		opts.AuthorID = v
	}

	ctx := context.Background()
	article, err := client.CreateArticle(ctx, portalID, opts)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(article)
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
