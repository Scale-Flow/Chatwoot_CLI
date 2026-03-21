# Security Policy

## Reporting a Vulnerability

Email scaleflowsolutions1@gmail.com with a description of the issue, steps to reproduce, and the affected version. Do not open a public GitHub issue for security problems.

You will receive acknowledgment within 48 hours. We aim to release a fix within 14 days of confirming the vulnerability.

## Scope

The following areas are in scope for security reports:

- **Credential handling**: API tokens stored in the OS keychain or read from environment variables. This includes token leakage through logs, error messages, or stdout output.
- **Token exposure**: Verbose or debug modes printing secrets to stderr. The CLI redacts tokens in log output, but report any case where redaction fails.
- **Command injection**: User-supplied input (profile names, filter values, resource identifiers) passed to shell commands or used in unsafe string interpolation.
- **Configuration file parsing**: Malicious YAML in `~/.config/chatwoot-cli/config.yaml` leading to code execution or file access outside the expected paths.

## Out of Scope

- Vulnerabilities in the Chatwoot server itself (report those to the [Chatwoot project](https://github.com/chatwoot/chatwoot/security)).
- Denial of service through large API responses.
