package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	innerauth "github.com/chatwoot/chatwoot-cli/internal/auth"
	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get authenticated user profile",
	RunE:  runProfileGet,
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update authenticated user profile",
	RunE:  runProfileUpdate,
}

func init() {
	profileUpdateCmd.Flags().String("name", "", "Display name")
	profileUpdateCmd.Flags().String("email", "", "Email address")
	profileUpdateCmd.Flags().String("availability", "", "Availability: online, offline, busy")
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileUpdateCmd)
}

// runtimeContext holds the resolved runtime configuration for a command.
type runtimeContext struct {
	ProfileName string
	BaseURL     string
	AccountID   int
}

// resolveContext resolves base URL, account ID, and profile from flags, env, and config.
func resolveContext(cmd *cobra.Command) (*runtimeContext, error) {
	flagProfile, _ := cmd.Root().PersistentFlags().GetString("profile")
	flagBaseURL, _ := cmd.Root().PersistentFlags().GetString("base-url")
	flagAccountID, _ := cmd.Root().PersistentFlags().GetInt("account-id")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")

	cfg, err := config.LoadFrom(cfgPath)
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && errors.Is(pathErr.Err, os.ErrNotExist) {
			cfg = &config.Config{}
		} else if _, statErr := os.Stat(cfgPath); errors.Is(statErr, os.ErrNotExist) {
			cfg = &config.Config{}
		} else {
			return nil, fmt.Errorf("config_error: %w", err)
		}
	}

	profileName, profile, err := cfg.ResolveProfile(flagProfile)
	if err != nil {
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

	return &runtimeContext{
		ProfileName: profileName,
		BaseURL:     resolved.BaseURL,
		AccountID:   resolved.AccountID,
	}, nil
}

// resolveAuth resolves credentials for the given profile and auth mode.
func resolveAuth(profileName string, mode credentials.AuthMode) (innerauth.TokenAuth, error) {
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
		return innerauth.ResolveApplication(resolver, profileName)
	case credentials.ModePlatform:
		return innerauth.ResolvePlatform(resolver, profileName)
	default:
		return innerauth.TokenAuth{}, fmt.Errorf("unknown auth mode: %s", mode)
	}
}

func prettyFromRoot(cmd *cobra.Command) bool {
	pretty, _ := cmd.Root().PersistentFlags().GetBool("pretty")
	return pretty
}

func writeError(cmd *cobra.Command, code, message string) error {
	resp := contract.Err(code, message)
	_ = contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
	return errors.New(message)
}

func runProfileGet(cmd *cobra.Command, args []string) error {
	rctx, err := resolveContext(cmd)
	if err != nil {
		return writeError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := resolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return writeError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	profile, err := client.GetProfile(context.Background())
	if err != nil {
		return writeError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func runProfileUpdate(cmd *cobra.Command, args []string) error {
	// Validate at least one flag is set
	nameChanged := cmd.Flags().Changed("name")
	emailChanged := cmd.Flags().Changed("email")
	availChanged := cmd.Flags().Changed("availability")
	if !nameChanged && !emailChanged && !availChanged {
		return fmt.Errorf("requires at least one of --name, --email, or --availability")
	}

	rctx, err := resolveContext(cmd)
	if err != nil {
		return writeError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := resolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return writeError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	opts := appapi.UpdateProfileOpts{}
	if nameChanged {
		v, _ := cmd.Flags().GetString("name")
		opts.Name = &v
	}
	if emailChanged {
		v, _ := cmd.Flags().GetString("email")
		opts.Email = &v
	}
	if availChanged {
		v, _ := cmd.Flags().GetString("availability")
		opts.Availability = &v
	}

	profile, err := client.UpdateProfile(context.Background(), opts)
	if err != nil {
		return writeError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(profile)
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}
