package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddAndGetPaginationFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddPaginationFlags(cmd)

	// Defaults
	pf := GetPaginationFlags(cmd)
	if pf.Page != 1 {
		t.Errorf("Page = %d, want 1", pf.Page)
	}
	if pf.PerPage != 25 {
		t.Errorf("PerPage = %d, want 25", pf.PerPage)
	}
	if pf.All {
		t.Error("All = true, want false")
	}
}

func TestPrettyFromRoot(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().Bool("pretty", false, "")
	child := &cobra.Command{Use: "child"}
	root.AddCommand(child)

	if Pretty(child) {
		t.Error("Pretty = true, want false")
	}
}
