# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-06-01

### Added

- Application API commands: conversations, contacts, messages, agents, inboxes, teams, labels, webhooks, reports, help center, custom attributes, custom filters, canned responses, and automation rules.
- Platform API commands: accounts, users, account-users, and agent-bots.
- Client API commands: contacts, conversations, and messages.
- Named profile support with YAML configuration at `~/.config/chatwoot-cli/config.yaml`.
- OS keychain credential storage via `go-keyring`, with environment variable fallback for CI environments.
- JSON envelope output on stdout for all commands, structured for machine consumption.
- Structured diagnostics on stderr via `log/slog`.
- Shell completion generators for bash, zsh, fish, and powershell.
- `--profile` flag to switch between named Chatwoot instances.
