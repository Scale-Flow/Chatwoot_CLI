package credentials

import "errors"

var ErrNotFound = errors.New("credential not found")

type AuthMode string

const (
	ModeApplication AuthMode = "application"
	ModePlatform    AuthMode = "platform"
)

type Store interface {
	Get(profile string, mode AuthMode) (string, error)
	Set(profile string, mode AuthMode, token string) error
	Delete(profile string, mode AuthMode) error
}
