package credentials

import (
	"errors"
	"os"
)

type EnvStore struct{}

func (s *EnvStore) Get(_ string, mode AuthMode) (string, error) {
	var envVar string
	switch mode {
	case ModeApplication:
		envVar = "CHATWOOT_ACCESS_TOKEN"
	case ModePlatform:
		envVar = "CHATWOOT_PLATFORM_TOKEN"
	default:
		return "", ErrNotFound
	}
	token := os.Getenv(envVar)
	if token == "" {
		return "", ErrNotFound
	}
	return token, nil
}

func (s *EnvStore) Set(_ string, _ AuthMode, _ string) error {
	return errors.New("cannot set credentials via environment store")
}

func (s *EnvStore) Delete(_ string, _ AuthMode) error {
	return errors.New("cannot delete credentials via environment store")
}
