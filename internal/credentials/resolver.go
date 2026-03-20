package credentials

type Source string

const (
	SourceEnv      Source = "environment"
	SourceKeychain Source = "keychain"
	SourceFile     Source = "file"
)

type Resolver struct {
	env      Store
	keychain Store
	file     Store
}

func NewResolver(env, keychain, file Store) *Resolver {
	return &Resolver{env: env, keychain: keychain, file: file}
}

func (r *Resolver) Get(profile string, mode AuthMode) (string, Source, error) {
	backends := []struct {
		store  Store
		source Source
	}{
		{r.env, SourceEnv},
		{r.keychain, SourceKeychain},
		{r.file, SourceFile},
	}
	for _, b := range backends {
		if b.store == nil {
			continue
		}
		token, err := b.store.Get(profile, mode)
		if err == ErrNotFound {
			continue
		}
		if err != nil {
			continue
		}
		return token, b.source, nil
	}
	return "", "", ErrNotFound
}
