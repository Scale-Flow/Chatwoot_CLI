package auth

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Store a credential for the active profile",
	RunE:  runSet,
}

func init() {
	setCmd.Flags().String("mode", "", "Auth mode: application or platform (required)")
	setCmd.Flags().String("token", "", "Token value (required)")
	_ = setCmd.MarkFlagRequired("mode")
	_ = setCmd.MarkFlagRequired("token")
	Cmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	modeStr, _ := cmd.Flags().GetString("mode")
	token, _ := cmd.Flags().GetString("token")

	mode := credentials.AuthMode(modeStr)
	if mode != credentials.ModeApplication && mode != credentials.ModePlatform {
		return fmt.Errorf("invalid mode %q: must be \"application\" or \"platform\"", modeStr)
	}

	profileName := resolveProfileNameForAuth(cmd)

	source, err := storeCredential(profileName, mode, token)
	if err != nil {
		resp := contract.Err(contract.ErrCodeAuth, err.Error())
		return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
	}

	resp := contract.Success(map[string]string{
		"profile": profileName,
		"mode":    string(mode),
		"source":  string(source),
	})
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func storeCredential(profile string, mode credentials.AuthMode, token string) (credentials.Source, error) {
	ks := credentials.NewKeychainStore()
	if err := ks.Set(profile, mode, token); err == nil {
		return credentials.SourceKeychain, nil
	} else {
		slog.Warn("keychain unavailable, using file-based credential storage", "error", err)
	}

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")

	fs := credentials.NewFileStore(credPath)
	if err := fs.Set(profile, mode, token); err != nil {
		return "", fmt.Errorf("store credential: %w", err)
	}
	return credentials.SourceFile, nil
}

// resolveProfileNameForAuth resolves just the profile name for auth commands.
// Does not validate BaseURL (auth commands don't need it).
func resolveProfileNameForAuth(cmd *cobra.Command) string {
	flagProfile, _ := cmd.Flags().GetString("profile")
	if flagProfile != "" {
		return flagProfile
	}
	if env := os.Getenv("CHATWOOT_PROFILE"); env != "" {
		return env
	}

	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		return "default"
	}
	if cfg.DefaultProfile != "" {
		return cfg.DefaultProfile
	}
	return "default"
}

func prettyFromRoot(cmd *cobra.Command) bool {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	return pretty
}
