# chatwoot-cli

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/Scale-Flow/Chatwoot_CLI?include_prereleases&sort=semver)](https://github.com/Scale-Flow/Chatwoot_CLI/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Scale-Flow/Chatwoot_CLI)](https://goreportcard.com/report/github.com/Scale-Flow/Chatwoot_CLI)

A machine-friendly CLI for the Chatwoot API. JSON on stdout, diagnostics on stderr, zero interactive prompts. Built for AI agents, scripts, and automation.

## Install

With `go install`:

```bash
go install github.com/chatwoot/chatwoot-cli/cmd/chatwoot@latest
```

From source (requires [Task](https://taskfile.dev)):

```bash
git clone https://github.com/chatwoot/chatwoot-cli.git
cd chatwoot-cli
task build
```

Or with plain Go:

```bash
go build -o chatwoot ./cmd/chatwoot/
```

## Quick start

Point the CLI at your Chatwoot instance and store your access token:

```bash
chatwoot auth set --base-url https://chatwoot.example.com --mode application --token <your-token>
```

Set a default account:

```bash
export CHATWOOT_ACCOUNT_ID=1
```

List conversations:

```bash
chatwoot application conversations list
```

Pipe through `jq`, feed into another tool, or parse the JSON envelope directly. Every response follows the same `{"ok": true, "data": ...}` shape.

## Command overview

The CLI mirrors the three Chatwoot API families. Each requires different credentials.

| Command | API | Auth |
|---|---|---|
| `chatwoot application` | Agent/admin (`/api/v1`, `/api/v2`) | User access token |
| `chatwoot platform` | Self-hosted admin (`/platform/api/v1`) | Platform app token |
| `chatwoot client` | Public end-user (`/public/api/v1`) | Inbox + contact identifiers |

**Application** covers the full operational surface: conversations, contacts, messages, agents, teams, inboxes, reports, webhooks, automation rules, canned responses, custom attributes, custom filters, labels, help center, integrations, audit logs, agent bots, and account settings.

**Platform** manages accounts, users, account-user mappings, and agent bots at the instance level.

**Client** exposes the public-facing subset: contacts, conversations, and messages.

Run `chatwoot <family> --help` to see available subcommands.

## Configuration

Config file location: `~/.config/chatwoot-cli/config.yaml`

```yaml
default_profile: work

profiles:
  work:
    base_url: https://chatwoot.example.com
    account_id: 1
  staging:
    base_url: https://staging.chatwoot.example.com
    account_id: 2
```

Select a profile per-command with `--profile staging`, or set it globally:

```bash
export CHATWOOT_PROFILE=staging
```

**Precedence order:** flags, then env vars, then config file, then defaults.

| Setting | Flag | Env var |
|---|---|---|
| Base URL | `--base-url` | `CHATWOOT_BASE_URL` |
| Account ID | `--account-id` | `CHATWOOT_ACCOUNT_ID` |
| Profile | `--profile` | `CHATWOOT_PROFILE` |

## Authentication

Store credentials in your OS keychain:

```bash
chatwoot auth set --mode application --token <user-access-token>
chatwoot auth set --mode platform --token <platform-app-token>
```

The CLI reads from the keychain at runtime. For CI environments where no keychain exists, set env vars instead:

```bash
export CHATWOOT_ACCESS_TOKEN=<user-access-token>     # application commands
export CHATWOOT_PLATFORM_TOKEN=<platform-app-token>   # platform commands
```

## JSON output

Every command writes a JSON envelope to stdout:

```json
{"ok": true, "data": {...}, "meta": {"pagination": {...}}}
```

```json
{"ok": false, "error": {"code": "not_found", "message": "Conversation not found", "detail": null}}
```

Pass `--pretty` for indented output:

```bash
chatwoot application conversations get --id 42 --pretty
```

Diagnostics and verbose logs go to stderr (`--verbose` to enable), so piping stdout to another program stays clean.

## Shell completion

Generate completions for your shell:

```bash
chatwoot completion bash > /etc/bash_completion.d/chatwoot
chatwoot completion zsh > "${fpath[1]}/_chatwoot"
chatwoot completion fish > ~/.config/fish/completions/chatwoot.fish
```

## Claude Code skill

This repo ships a [Claude Code skill](https://docs.anthropic.com/en/docs/claude-code/skills) that teaches Claude the CLI's command syntax, flag names, JSON envelope handling, and common workflows. With the skill loaded, Claude uses the right commands on the first try instead of discovering them through `--help`.

### Install the skill

Copy the skill into your Claude Code skills directory:

```bash
mkdir -p ~/.claude/skills/chatwoot-cli
cp .claude/skills/chatwoot-cli/SKILL.md ~/.claude/skills/chatwoot-cli/SKILL.md
```

The skill triggers automatically when you mention Chatwoot in a CLI or automation context.

### What the skill provides

- Exact flag names for every command, verified against the binary
- JSON envelope parsing patterns (`ok`/`data`/`error`)
- Auth setup for all three API families
- Pagination with `--page` and `--all`
- Common multi-step workflows: conversation reply, contact lookup, report generation

## Documentation

- [Command reference](docs/commands.md)
- [Configuration guide](docs/configuration.md)
- [Contributing](CONTRIBUTING.md)

## License

MIT. Copyright Scaleflow.
