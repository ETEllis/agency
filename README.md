# TeamCode

⌬ Terminal-based AI coding team assistant with multi-agent collaboration support.

## Overview

TeamCode is a fork of OpenCode focused on first-class multi-agent collaboration. It keeps the terminal-native workflow, model selection, and coding tools from OpenCode, while adding persistent team coordination and formal subagent execution. It enables:

- **Team Context**: Shared charter, roles, goals, and working agreements
- **Task Board**: Shared backlog, in-progress, blocked, and done state
- **Handoff Protocol**: Explicit task transitions between agents
- **Private Inboxes**: Each agent has their own message inbox
- **Direct Chatter**: Teammates can message specific peers
- **Global Broadcasts**: Team-wide updates and coordination signals
- **Persistent Teammates**: Named worker sessions with role identity
- **Subagents**: Bounded worker sessions for focused delegated execution
- **Nested Delegation**: Teammates can spawn their own subagents
- **Leader-Led Execution**: Teams can be bootstrapped around a lead who coordinates delegation, broadcasts, and task board flow

## Architecture

```
TeamCode (Go)
    │
    ├── internal/team/             # Shared team state and messaging
    │       ├── team_context.json
    │       ├── task_board.json
    │       ├── handoffs.json
    │       ├── members.json
    │       └── inboxes/<agent>.json
    │
    └── internal/orchestration/    # Persistent teammate/subagent runtime
            └── worker manager + child sessions
```

Team state is now native Go and file-backed under `.teamcode/teams` or the configured TeamCode data directory. The old Python bridge is no longer part of the runtime.

## Installation

### Using the Install Script

```bash
# Install the latest version
curl -fsSL https://raw.githubusercontent.com/ETEllis/teamcode/refs/heads/main/install | bash

# Install a specific version
curl -fsSL https://raw.githubusercontent.com/ETEllis/teamcode/refs/heads/main/install | VERSION=0.1.0 bash
```

### Using Homebrew

```bash
brew install --cask ETEllis/tap/teamcode
```

### Using Go

```bash
go install github.com/ETEllis/teamcode@latest
```

## Building From Source

```bash
go build -o teamcode .
```

## Usage

```bash
# Interactive mode
./teamcode

# Non-interactive mode
./teamcode -p "Explain this codebase"

# With debug logging
./teamcode -d
```

## Configuration

TeamCode looks for configuration in these locations:

- `$HOME/.teamcode.json`
- `$XDG_CONFIG_HOME/teamcode/.teamcode.json`
- `./.teamcode.json`

Legacy OpenCode config files are still read as fallbacks:

- `$HOME/.opencode.json`
- `$XDG_CONFIG_HOME/opencode/.opencode.json`
- `./.opencode.json`

You can also configure providers with environment variables such as `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, `GITHUB_TOKEN`, `OPENROUTER_API_KEY`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `AZURE_OPENAI_ENDPOINT`, and `AZURE_OPENAI_API_KEY`.

## Collaboration Tools

The main coder agent now has first-class collaboration tools:

- `team_bootstrap` - Create a leader-led team with default specialist roles and optional spawned teammates
- `team_create_context` - Create team with charter and roles
- `team_add_role` - Add role definitions
- `team_assign_role` - Assign agents to roles
- `task_create` - Add tasks to the board
- `task_move` - Move tasks between columns
- `handoff_create` - Create task handoffs
- `handoff_accept` - Accept pending handoffs
- `inbox_read` - Read your messages
- `team_message_send` - Send a direct teammate message
- `team_broadcast` - Send a team-wide broadcast
- `team_status` - Inspect team, task, member, handoff, and worker state
- `teammate_spawn` / `teammate_wait` - Launch and monitor persistent teammates
- `subagent_spawn` / `subagent_wait` - Launch and monitor bounded subagents

For most coordinated workflows, TeamCode should start with `team_bootstrap` and a leader-led pattern instead of manually stitching roles together one tool at a time.

## Memory and Commands

- Team memory files: `TeamCode.md`, `teamcode.md`
- Legacy memory files still supported: `OpenCode.md`, `opencode.md`
- User custom commands: `$HOME/.teamcode/commands/`
- Project custom commands: `<PROJECT>/.teamcode/commands/`

## Compatibility

TeamCode prefers `.teamcode`, `TeamCode.md`, `teamcode.db`, and `teamcode` theme names, but still reads legacy OpenCode config and memory files so existing installs can migrate without breaking.

## License

MIT
