package cli

import (
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/spf13/cobra"
)

// ResolveProfileName resolves just the profile name (lightweight path for auth commands).
func ResolveProfileName(cmd *cobra.Command) string {
	flagProfile, _ := cmd.Flags().GetString("profile")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		cfg = &config.Config{}
	}

	name, _, err := cfg.ResolveProfile(flagProfile)
	if err != nil {
		// Fall back to flag, env, or "default"
		name = flagProfile
		if name == "" {
			name = os.Getenv("CHATWOOT_PROFILE")
		}
		if name == "" {
			name = "default"
		}
	}
	return name
}
