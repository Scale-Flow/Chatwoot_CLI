package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Profile holds connection settings for a single Chatwoot instance.
type Profile struct {
	BaseURL   string `mapstructure:"base_url"`
	AccountID int    `mapstructure:"account_id"`
}

// Config is the top-level configuration, containing named profiles.
type Config struct {
	DefaultProfile string             `mapstructure:"default_profile"`
	Profiles       map[string]Profile `mapstructure:"profiles"`
}

// LoadFrom reads and parses a YAML config file at path.
func LoadFrom(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}

// ResolveProfile selects a profile by precedence: flag > CHATWOOT_PROFILE env > default_profile > "default".
func (c *Config) ResolveProfile(flagProfile string) (string, Profile, error) {
	name := flagProfile
	if name == "" {
		name = os.Getenv("CHATWOOT_PROFILE")
	}
	if name == "" {
		name = c.DefaultProfile
	}
	if name == "" {
		name = "default"
	}
	profile, ok := c.Profiles[name]
	if !ok {
		return "", Profile{}, fmt.Errorf("profile %q not found in config", name)
	}
	return name, profile, nil
}

// ResolveOverrides applies flag and env var overrides to a profile.
// Precedence: flag > env var > profile value.
func ResolveOverrides(p Profile, flagBaseURL string, flagAccountID int) Profile {
	if flagBaseURL != "" {
		p.BaseURL = flagBaseURL
	} else if env := os.Getenv("CHATWOOT_BASE_URL"); env != "" {
		p.BaseURL = env
	}
	if flagAccountID != 0 {
		p.AccountID = flagAccountID
	} else if env := os.Getenv("CHATWOOT_ACCOUNT_ID"); env != "" {
		if id, err := strconv.Atoi(env); err == nil {
			p.AccountID = id
		}
	}
	return p
}
