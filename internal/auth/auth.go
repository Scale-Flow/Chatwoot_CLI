package auth

import (
	"fmt"

	"github.com/chatwoot/chatwoot-cli/internal/credentials"
)

type TokenAuth struct {
	Token      string
	HeaderName string
	Source     credentials.Source
}

type ClientAuth struct {
	InboxIdentifier   string
	ContactIdentifier string
}

func ResolveApplication(resolver *credentials.Resolver, profile string) (TokenAuth, error) {
	token, source, err := resolver.Get(profile, credentials.ModeApplication)
	if err != nil {
		return TokenAuth{}, fmt.Errorf("no application credentials for profile %q: %w", profile, err)
	}
	return TokenAuth{Token: token, HeaderName: "api_access_token", Source: source}, nil
}

func ResolvePlatform(resolver *credentials.Resolver, profile string) (TokenAuth, error) {
	token, source, err := resolver.Get(profile, credentials.ModePlatform)
	if err != nil {
		return TokenAuth{}, fmt.Errorf("no platform credentials for profile %q: %w", profile, err)
	}
	return TokenAuth{Token: token, HeaderName: "api_access_token", Source: source}, nil
}
