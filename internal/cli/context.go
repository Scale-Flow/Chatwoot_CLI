package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/auth"
	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

// RuntimeContext holds the resolved runtime configuration for a command.
type RuntimeContext struct {
	ProfileName string
	BaseURL     string
	AccountID   int
}

// ResolveContext resolves the full runtime context from flags, env, and config.
func ResolveContext(cmd *cobra.Command) (*RuntimeContext, error) {
	flagProfile, _ := cmd.Flags().GetString("profile")
	flagBaseURL, _ := cmd.Flags().GetString("base-url")
	flagAccountID, _ := cmd.Flags().GetInt("account-id")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	return resolveContextFromPath(cfgPath, flagProfile, flagBaseURL, flagAccountID)
}

// resolveContextFromPath is the testable core of ResolveContext.
func resolveContextFromPath(cfgPath, flagProfile, flagBaseURL string, flagAccountID int) (*RuntimeContext, error) {
	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		// Config file missing is OK — proceed with empty config.
		// Viper wraps the underlying error, so check with errors.As.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			cfg = &config.Config{}
		} else if _, statErr := os.Stat(cfgPath); errors.Is(statErr, os.ErrNotExist) {
			// Fallback: file genuinely doesn't exist
			cfg = &config.Config{}
		} else {
			return nil, fmt.Errorf("config_error: %w", err)
		}
	}

	profileName, profile, err := cfg.ResolveProfile(flagProfile)
	if err != nil {
		// If config is empty and profile not found, use empty profile
		if cfg.Profiles == nil || len(cfg.Profiles) == 0 {
			profileName = flagProfile
			if profileName == "" {
				profileName = os.Getenv("CHATWOOT_PROFILE")
			}
			if profileName == "" {
				profileName = "default"
			}
			profile = config.Profile{}
		} else {
			return nil, fmt.Errorf("config_error: %w", err)
		}
	}

	resolved := config.ResolveOverrides(profile, flagBaseURL, flagAccountID)

	if resolved.BaseURL == "" {
		return nil, fmt.Errorf("no base URL configured — set base_url in profile or use --base-url flag")
	}

	return &RuntimeContext{
		ProfileName: profileName,
		BaseURL:     resolved.BaseURL,
		AccountID:   resolved.AccountID,
	}, nil
}

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

// ResolveAuth resolves credentials for the given profile and auth mode.
func ResolveAuth(profileName string, mode credentials.AuthMode) (auth.TokenAuth, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")

	resolver := credentials.NewResolver(
		&credentials.EnvStore{},
		credentials.NewKeychainStore(),
		credentials.NewFileStore(credPath),
	)

	switch mode {
	case credentials.ModeApplication:
		return auth.ResolveApplication(resolver, profileName)
	case credentials.ModePlatform:
		return auth.ResolvePlatform(resolver, profileName)
	default:
		return auth.TokenAuth{}, fmt.Errorf("unknown auth mode: %s", mode)
	}
}

// WriteError writes an error envelope to stdout and returns an error for Cobra.
func WriteError(cmd *cobra.Command, code, message string) error {
	resp := contract.Err(code, message)
	_ = contract.Write(cmd.OutOrStdout(), resp, prettyFlag)
	return errors.New(message)
}
