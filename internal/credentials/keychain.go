package credentials

import (
	"errors"
	"fmt"

	keyring "github.com/zalando/go-keyring"
)

const keychainService = "chatwoot-cli"

type KeychainStore struct{}

func NewKeychainStore() *KeychainStore {
	return &KeychainStore{}
}

func (s *KeychainStore) Get(profile string, mode AuthMode) (string, error) {
	user := fmt.Sprintf("%s/%s", profile, mode)
	token, err := keyring.Get(keychainService, user)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("keychain get: %w", err)
	}
	return token, nil
}

func (s *KeychainStore) Set(profile string, mode AuthMode, token string) error {
	user := fmt.Sprintf("%s/%s", profile, mode)
	if err := keyring.Set(keychainService, user, token); err != nil {
		return fmt.Errorf("keychain set: %w", err)
	}
	return nil
}

func (s *KeychainStore) Delete(profile string, mode AuthMode) error {
	user := fmt.Sprintf("%s/%s", profile, mode)
	if err := keyring.Delete(keychainService, user); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("keychain delete: %w", err)
	}
	return nil
}
