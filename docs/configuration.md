# Configuration Guide

## Config File

The CLI reads configuration from `~/.config/chatwoot-cli/config.yaml`.

A minimal config file:

```yaml
default_profile: production

profiles:
  production:
    base_url: https://app.chatwoot.com
    account_id: 1
```

Every setting under a profile maps to a global flag. The file uses standard YAML.

## Profiles

Profiles let you store connection details for separate Chatwoot instances. Each profile has a name, a base URL, and an account ID.

```yaml
default_profile: production

profiles:
  production:
    base_url: https://app.chatwoot.com
    account_id: 1
  staging:
    base_url: https://staging.example.com
    account_id: 2
```

Switch profiles per-command with the `--profile` flag:

```bash
chatwoot application conversations list --profile staging
```

The CLI resolves which profile to use in this order:

1. `--profile` flag
2. `CHATWOOT_PROFILE` environment variable
3. `default_profile` field in the config file
4. A profile named `default`

Credentials bind to the active profile. When you run `chatwoot auth set`, the token stores under the profile name currently in scope.

## Authentication

Two auth modes exist, matching the two authenticated API families:

- **application** targets the agent/admin API (`/api/v1`, `/api/v2`) with a user access token.
- **platform** targets the self-hosted admin API (`/platform/api/v1`) with a platform app token.

Store a token in the OS keychain:

```bash
chatwoot auth set --mode application --token sk-abc123
chatwoot auth set --mode platform --token pt-xyz789
```

Check what the CLI can see:

```bash
chatwoot auth status
```

Remove credentials for the active profile:

```bash
chatwoot auth clear
```

The CLI checks the OS keychain first, then falls back to environment variables (`CHATWOOT_ACCESS_TOKEN` or `CHATWOOT_PLATFORM_TOKEN`). Keychain storage uses the profile name as the lookup key.

The Client API (`/public/api/v1`) skips token auth entirely. Pass `--inbox-id` and `--contact-id` flags on each command instead.

## Environment Variables

| Variable | Purpose |
|---|---|
| `CHATWOOT_BASE_URL` | Override the base URL for the active profile |
| `CHATWOOT_ACCOUNT_ID` | Override the account ID for the active profile |
| `CHATWOOT_PROFILE` | Select a named profile |
| `CHATWOOT_ACCESS_TOKEN` | Application API token (keychain fallback) |
| `CHATWOOT_PLATFORM_TOKEN` | Platform API token (keychain fallback) |

## Precedence Rules

Settings resolve left to right, with the leftmost source winning:

| Setting | Flag | Env Var | Config File | Default |
|---|---|---|---|---|
| Base URL | `--base-url` | `CHATWOOT_BASE_URL` | `profiles.<name>.base_url` | none |
| Account ID | `--account-id` | `CHATWOOT_ACCOUNT_ID` | `profiles.<name>.account_id` | none |
| Profile | `--profile` | `CHATWOOT_PROFILE` | `default_profile` | `default` |
| JSON formatting | `--pretty` | -- | -- | compact |
| Diagnostic logs | `--verbose` | -- | -- | off |

Flags beat environment variables. Environment variables beat config file values. Config file values beat built-in defaults.

Credentials follow their own chain: OS keychain first, then `CHATWOOT_ACCESS_TOKEN` / `CHATWOOT_PLATFORM_TOKEN`.

## CI and Headless Usage

CI runners and containers lack a keychain. Set tokens through environment variables:

```bash
export CHATWOOT_BASE_URL=https://app.chatwoot.com
export CHATWOOT_ACCOUNT_ID=1
export CHATWOOT_ACCESS_TOKEN=sk-abc123
```

No config file or `auth set` step needed. The CLI picks up all three values from the environment and runs without interactive setup.

For pipelines that target multiple instances, set `CHATWOOT_PROFILE` alongside a config file, or override `CHATWOOT_BASE_URL` and `CHATWOOT_ACCOUNT_ID` per step.

All output goes to stdout as JSON. Diagnostic logs (`--verbose`) go to stderr. This separation keeps pipeline parsing clean regardless of log level.
