package cmdutil

import "github.com/spf13/cobra"

// PaginationFlags holds pagination flag values.
type PaginationFlags struct {
	Page    int
	PerPage int
	All     bool
}

// AddPaginationFlags registers --page, --per-page, and --all flags on cmd.
func AddPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 1, "Page number")
	cmd.Flags().Int("per-page", 25, "Items per page")
	cmd.Flags().Bool("all", false, "Fetch all pages")
}

// GetPaginationFlags reads pagination flag values from cmd.
func GetPaginationFlags(cmd *cobra.Command) PaginationFlags {
	page, _ := cmd.Flags().GetInt("page")
	perPage, _ := cmd.Flags().GetInt("per-page")
	all, _ := cmd.Flags().GetBool("all")
	return PaginationFlags{Page: page, PerPage: perPage, All: all}
}

// Pretty reads the --pretty flag from the root command.
func Pretty(cmd *cobra.Command) bool {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	return pretty
}
