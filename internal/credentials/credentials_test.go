package credentials

import "testing"

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
