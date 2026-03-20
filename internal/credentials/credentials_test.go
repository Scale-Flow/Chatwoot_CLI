package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvStoreGetApplication(t *testing.T) {
	t.Setenv("CHATWOOT_ACCESS_TOKEN", "sk-test-token")
	store := &EnvStore{}

	token, err := store.Get("any-profile", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "sk-test-token" {
		t.Errorf("token = %q, want %q", token, "sk-test-token")
	}
}

func TestEnvStoreGetPlatform(t *testing.T) {
	t.Setenv("CHATWOOT_PLATFORM_TOKEN", "pk-test-token")
	store := &EnvStore{}

	token, err := store.Get("any-profile", ModePlatform)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "pk-test-token" {
		t.Errorf("token = %q, want %q", token, "pk-test-token")
	}
}

func TestEnvStoreGetNotSet(t *testing.T) {
	store := &EnvStore{}
	_, err := store.Get("any", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestFileStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	err := store.Set("work", ModeApplication, "sk-secret")
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %o, want 600", info.Mode().Perm())
	}

	token, err := store.Get("work", ModeApplication)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if token != "sk-secret" {
		t.Errorf("token = %q, want %q", token, "sk-secret")
	}
}

func TestFileStoreGetNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	_, err := store.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestFileStoreDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")
	store := NewFileStore(path)

	_ = store.Set("work", ModeApplication, "sk-secret")
	err := store.Delete("work", ModeApplication)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = store.Get("work", ModeApplication)
	if err != ErrNotFound {
		t.Errorf("after delete: err = %v, want ErrNotFound", err)
	}
}

func TestFileStoreRejectsWidePerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credentials.yaml")

	err := os.WriteFile(path, []byte("profiles: {}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	store := NewFileStore(path)
	_, err = store.Get("work", ModeApplication)
	if err == nil {
		t.Fatal("expected error for wide permissions")
	}
}
