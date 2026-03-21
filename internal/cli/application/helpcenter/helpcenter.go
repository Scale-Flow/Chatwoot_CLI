package helpcenter

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{Use: "help-center", Short: "Manage help center"}

var portalsCmd = &cobra.Command{Use: "portals", Short: "Manage portals"}
var articlesCmd = &cobra.Command{Use: "articles", Short: "Manage articles"}
var categoriesCmd = &cobra.Command{Use: "categories", Short: "Manage categories"}

func init() {
	Cmd.AddCommand(portalsCmd)
	Cmd.AddCommand(articlesCmd)
	Cmd.AddCommand(categoriesCmd)
}
