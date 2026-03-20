package auth

import (
	"os"
	"path/filepath"

	"github.com/chatwoot/chatwoot-cli/internal/config"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show credential status for the active profile",
	RunE:  runStatus,
}

func init() {
	Cmd.AddCommand(statusCmd)
}

type credentialStatus struct {
	Status string `json:"status"`
	Source string `json:"source,omitempty"`
}

type statusData struct {
	Profile     string                      `json:"profile"`
	BaseURL     string                      `json:"base_url,omitempty"`
	AccountID   int                         `json:"account_id,omitempty"`
	Credentials map[string]credentialStatus `json:"credentials"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	profileName := resolveProfileNameForAuth(cmd)

	// Try to load config for display (non-fatal if missing)
	var baseURL string
	var accountID int
	cfgDir, _ := os.UserConfigDir()
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(cfgDir, "chatwoot-cli", "config.yaml")
	if cfg, err := config.LoadFrom(cfgPath); err == nil {
		if _, profile, err := cfg.ResolveProfile(profileName); err == nil {
			baseURL = profile.BaseURL
			accountID = profile.AccountID
		}
	}

	// Override with env/flags
	resolved := config.ResolveOverrides(config.Profile{BaseURL: baseURL, AccountID: accountID}, "", 0)
	baseURL = resolved.BaseURL
	accountID = resolved.AccountID

	// Probe credentials
	credPath := filepath.Join(cfgDir, "chatwoot-cli", "credentials.yaml")
	resolver := credentials.NewResolver(
		&credentials.EnvStore{},
		credentials.NewKeychainStore(),
		credentials.NewFileStore(credPath),
	)

	creds := map[string]credentialStatus{
		"application": probeCredential(resolver, profileName, credentials.ModeApplication),
		"platform":    probeCredential(resolver, profileName, credentials.ModePlatform),
	}

	result := statusData{
		Profile:     profileName,
		BaseURL:     baseURL,
		AccountID:   accountID,
		Credentials: creds,
	}

	resp := contract.Success(result)
	return contract.Write(cmd.OutOrStdout(), resp, prettyFromRoot(cmd))
}

func probeCredential(resolver *credentials.Resolver, profile string, mode credentials.AuthMode) credentialStatus {
	_, source, err := resolver.Get(profile, mode)
	if err != nil {
		return credentialStatus{Status: "not_configured"}
	}
	return credentialStatus{Status: "configured", Source: string(source)}
}
