package auth

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove credentials for the active profile",
	RunE:  runClear,
	PostRunE: func(cmd *cobra.Command, args []string) error {
		// Reset flags to defaults so shared command instances work correctly in tests.
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			f.Changed = false
			_ = f.Value.Set(f.DefValue)
		})
		return nil
	},
}

func init() {
	clearCmd.Flags().String("mode", "", "Auth mode to clear: application or platform")
	clearCmd.Flags().Bool("all", false, "Clear all credential modes")
	Cmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	modeStr, _ := cmd.Flags().GetString("mode")
	all, _ := cmd.Flags().GetBool("all")

	if modeStr == "" && !all {
		return fmt.Errorf("requires --mode or --all flag")
	}

	profileName := resolveProfileNameForAuth(cmd)

	var modes []credentials.AuthMode
	if all {
		modes = []credentials.AuthMode{credentials.ModeApplication, credentials.ModePlatform}
	} else {
		mode := credentials.AuthMode(modeStr)
		if mode != credentials.ModeApplication && mode != credentials.ModePlatform {
			return fmt.Errorf("invalid mode %q: must be \"application\" or \"platform\"", modeStr)
		}
		modes = []credentials.AuthMode{mode}
	}

	var cleared []string
	for _, mode := range modes {
		if err := deleteCredential(profileName, mode); err != nil {
			slog.Warn("failed to delete credential", "mode", mode, "error", err)
		}
		cleared = append(cleared, string(mode))
	}

	resp := contract.Success(map[string]any{
		"profile": profileName,
		"cleared": cleared,
	})
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func deleteCredential(profile string, mode credentials.AuthMode) error {
	ks := credentials.NewKeychainStore()
	if err := ks.Delete(profile, mode); err == nil {
		return nil
	}

	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")
	fs := credentials.NewFileStore(credPath)
	return fs.Delete(profile, mode)
}
