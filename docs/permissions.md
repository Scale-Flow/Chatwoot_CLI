# Chatwoot CLI Permissions & Access Control

## Overview

The Chatwoot CLI is designed for use by AI agents and automation. Since Chatwoot API tokens are not independently scoped (a token inherits the full permissions of its user's role), the permission model relies on a combination of **server-side RBAC** and **client-side guards**.

## Chatwoot's Permission Hierarchy

Chatwoot has three levels of access, managed in two separate interfaces:

| Level | Role | Scope | Managed In |
|-------|------|-------|------------|
| Installation | Super Admin | Entire Chatwoot server — all accounts, system config | `/super_admin` console |
| Account | Administrator | Full control within one account — settings, inboxes, agents, billing | Settings → Agents |
| Account | Agent | Conversations, contacts, read-only access to config | Settings → Agents |
| Account | Custom Role (Enterprise) | Granular permissions within one account | Settings → Custom Roles |

**Important:** The Super Admin console (`/super_admin`) only manages installation-level super admins. Account-level roles (Administrator, Agent, Custom) are managed in the regular Chatwoot dashboard under **Settings → Agents**.

### Enterprise Custom Role Permissions

On Chatwoot Enterprise, custom roles provide granular control:

- `conversation_manage` — view, update, assign, close conversations
- `conversation_unassigned_manage` — handle unassigned conversations only
- `contact_manage` — create, update, delete, merge contacts
- `inbox_manage` — configure inboxes
- `team_manage` — create and manage teams
- `report_manage` — access analytics and exports

## API Authentication Models

The CLI supports three API families, each with a different auth model:

| API Family | Token Type | Header | Availability |
|-----------|-----------|--------|-------------|
| Application (`/api/v1`, `/api/v2`) | User Access Token | `api_access_token` | Cloud + Self-hosted |
| Platform (`/platform/api/v1`) | Platform App Token | `api_access_token` | Self-hosted only |
| Client (`/public/api/v1`) | Inbox + Contact identifiers | Path-based | Cloud + Self-hosted |

### How to Obtain Tokens

- **User Access Token:** Log in → click avatar (bottom-left) → Profile Settings → scroll to bottom. Each user gets one auto-generated token. It does not expire but can be reset.
- **Platform App Token:** Super Admin console → Platform Apps → New Platform App → Access Tokens tab. Self-hosted only.
- **Client Identifiers:** `inbox_identifier` from Settings → Inboxes → Configuration tab. `contact_identifier` from `POST /public/api/v1/inboxes/{inbox_identifier}/contacts`.

## Server-Side Permission Matrix (Verified)

Tested against Chatwoot v4.x. The Agent role provides meaningful server-side restrictions.

### Read Operations

| Endpoint | Agent | Admin |
|----------|-------|-------|
| `GET /profile` | ✅ 200 | ✅ 200 |
| `GET /conversations` | ✅ 200 | ✅ 200 |
| `GET /contacts` | ✅ 200 | ✅ 200 |
| `GET /agents` | ✅ 200 | ✅ 200 |
| `GET /teams` | ✅ 200 | ✅ 200 |
| `GET /labels` | ✅ 200 | ✅ 200 |
| `GET /custom_attribute_definitions` | ✅ 200 | ✅ 200 |
| `GET /custom_filters` | ✅ 200 | ✅ 200 |
| `GET /canned_responses` | ✅ 200 | ✅ 200 |
| `GET /integrations/apps` | ✅ 200 | ✅ 200 |
| `GET /agent_bots` | ✅ 200 | ✅ 200 |
| `GET /webhooks` | ❌ 401 | ✅ 200 |
| `GET /automation_rules` | ❌ 401 | ✅ 200 |

### Write Operations

| Endpoint | Agent | Admin |
|----------|-------|-------|
| `POST /labels` | ❌ 401 | ✅ 200 |
| `POST /webhooks` | ❌ 401 | ✅ 200 |
| `POST /automation_rules` | ❌ 401 | ✅ 200 |
| `DELETE /agents/:id` | ❌ 401 | ✅ (expected) |

### Summary

**Agent role is blocked from:**
- Webhooks (read and write)
- Automation rules (read and write)
- Label management (create/update/delete)
- Agent management (delete)
- Account-level configuration

**Agent role can:**
- Read conversations, contacts, teams, inboxes, custom attributes, canned responses, integrations, agent bots
- Work with conversations (send messages, assign, update status)
- Update contacts

## CLI Permission Design

### Layered Defense Model

Since Chatwoot tokens cannot be independently scoped, the CLI uses three layers:

```
┌─────────────────────────────────────────────────┐
│  Layer 1: Server-Side RBAC (primary boundary)   │
│  Chatwoot enforces role permissions with 401s.   │
│  Create least-privilege users per use case.      │
├─────────────────────────────────────────────────┤
│  Layer 2: CLI-Side Guards                        │
│  --read-only flag prevents all mutations.        │
│  Profile config provides advisory scoping.       │
├─────────────────────────────────────────────────┤
│  Layer 3: Skill/Orchestrator                     │
│  The calling skill/agent only describes          │
│  commands relevant to the task.                  │
└─────────────────────────────────────────────────┘
```

### Layer 1: Server-Side RBAC (Enforcement)

The primary security boundary. Create dedicated Chatwoot users per use case:

| Use Case | Recommended Role | Rationale |
|----------|-----------------|-----------|
| Support agent automation | Agent | Can read/write conversations and contacts, blocked from admin ops |
| Read-only analytics | Agent + CLI `--read-only` | Agent role allows reads; CLI blocks mutations |
| Full account automation | Administrator | Use only when admin operations are genuinely required |
| Provisioning/multi-tenant | Platform App token | Isolated to objects created by that platform app |
| End-user messaging | Client API identifiers | Natural data isolation — can only access own conversations |

### Layer 2: CLI-Side Guards (Advisory + Mutation Prevention)

**Global `--read-only` flag:**

```bash
chatwoot --read-only conversations list          # allowed
chatwoot --read-only conversations assign --id 5  # blocked by CLI before API call
```

**Profile-level configuration:**

```yaml
profiles:
  support-bot:
    base_url: https://support.adminflow.ca
    token: "@keyring:support-bot"
    read_only: false

  analytics:
    base_url: https://support.adminflow.ca
    token: "@keyring:analytics"
    read_only: true   # CLI blocks all POST/PUT/PATCH/DELETE
```

**Important:** Profile-based scopes are advisory, not a security boundary. An agent with shell access can bypass config-file restrictions. The server-side role is the real enforcement.

### Layer 3: Skill/Orchestrator (Tool Surface Restriction)

When the CLI is used as a skill for AI agents, the skill prompt controls which commands the agent knows about. This is **not a security boundary** — it relies on the agent following instructions — but reduces the likelihood of unintended operations.

### Command Schema Introspection

Every CLI command declares its permission requirements for orchestrator use:

```bash
chatwoot schema                                    # full command schema
chatwoot schema --command conversations.list       # single command metadata
```

Output includes:

```json
{
  "command": "conversations.assign",
  "method": "POST",
  "mutation": true,
  "minimum_role": "agent",
  "api_family": "application"
}
```

This allows orchestrators and skill authors to programmatically determine which commands are safe for a given use case.

## Security Considerations for AI Agent Use

### The Self-Escalation Problem

AI agents with shell access can potentially modify config files to grant themselves more permissions. Mitigations:

1. **Server-side role is the real boundary.** Even if an agent modifies its CLI config, the Chatwoot API still enforces the user's role. A token from an Agent-role user will always get 401 on admin endpoints regardless of CLI config.

2. **Tokens in OS keychain.** Stored via `go-keyring`, not in config files. The agent cannot extract a higher-privilege token from the keychain without OS-level authentication.

3. **Token redaction.** The CLI never outputs tokens in JSON responses, error messages, or verbose logs. An agent cannot exfiltrate its own token to use via raw HTTP.

4. **Separate users per trust level.** Never share an Administrator token with an agent that only needs Agent-level access.

### Recommended Setup for AI Agent Use

1. Create a dedicated Chatwoot user with the **Agent** role
2. Store its token in the OS keychain under a named profile
3. Invoke the CLI with `--profile <name>` (and `--read-only` if appropriate)
4. The skill/orchestrator describes only the relevant subset of commands
5. Server-side RBAC enforces the actual permission boundary

### Token Hygiene

- Rotate tokens after any suspected exposure (Profile Settings → Reset)
- Never share tokens in chat, logs, or issue trackers
- Use environment variables (`CHATWOOT_ACCESS_TOKEN`) for CI/CD, not command-line flags (which appear in process lists)
- Platform App tokens should be treated with the same care as database credentials
