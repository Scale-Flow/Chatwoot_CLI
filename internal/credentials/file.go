package credentials

import (
	"fmt"
	"os"
	"sync"

	"go.yaml.in/yaml/v3"
)

type fileData struct {
	Profiles map[string]fileProfile `yaml:"profiles"`
}

type fileProfile struct {
	ApplicationToken string `yaml:"application_token,omitempty"`
	PlatformToken    string `yaml:"platform_token,omitempty"`
}

type FileStore struct {
	path string
	mu   sync.Mutex
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Get(profile string, mode AuthMode) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return "", err
	}

	p, ok := data.Profiles[profile]
	if !ok {
		return "", ErrNotFound
	}

	var token string
	switch mode {
	case ModeApplication:
		token = p.ApplicationToken
	case ModePlatform:
		token = p.PlatformToken
	}
	if token == "" {
		return "", ErrNotFound
	}
	return token, nil
}

func (s *FileStore) Set(profile string, mode AuthMode, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil && !os.IsNotExist(err) && err != ErrNotFound {
		return err
	}
	if data.Profiles == nil {
		data.Profiles = make(map[string]fileProfile)
	}

	p := data.Profiles[profile]
	switch mode {
	case ModeApplication:
		p.ApplicationToken = token
	case ModePlatform:
		p.PlatformToken = token
	}
	data.Profiles[profile] = p

	return s.save(data)
}

func (s *FileStore) Delete(profile string, mode AuthMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	p, ok := data.Profiles[profile]
	if !ok {
		return nil
	}

	switch mode {
	case ModeApplication:
		p.ApplicationToken = ""
	case ModePlatform:
		p.PlatformToken = ""
	}
	data.Profiles[profile] = p

	return s.save(data)
}

func (s *FileStore) load() (fileData, error) {
	var data fileData

	info, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		return data, ErrNotFound
	}
	if err != nil {
		return data, fmt.Errorf("stat credentials file: %w", err)
	}

	if info.Mode().Perm()&0077 != 0 {
		return data, fmt.Errorf("credentials file %s has permissions %o, want 0600 or stricter", s.path, info.Mode().Perm())
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return data, fmt.Errorf("read credentials file: %w", err)
	}

	if err := yaml.Unmarshal(raw, &data); err != nil {
		return data, fmt.Errorf("unmarshal credentials: %w", err)
	}
	return data, nil
}

func (s *FileStore) save(data fileData) error {
	raw, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal credentials: %w", err)
	}
	return os.WriteFile(s.path, raw, 0600)
}
