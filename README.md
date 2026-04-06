# PAI Universal - Personal AI Infrastructure

**Version:** 0.1.0-alpha  
**Status:** Phase 1 - Infrastructure Port

## Overview

PAI Universal (Personal AI Infrastructure) is a goal-oriented AI assistant framework that works across multiple AI tools. Unlike stateless AI assistants, PAI knows your goals, remembers your preferences, and continuously improves its understanding of how you work.

### Supported AI Tools

| Tool | Status | Notes |
|------|--------|-------|
| opencode | ✓ Primary | Native skill system |
| Claude Code | Legacy | Original PAI |
| Cursor | Planned | Future support |
| Codex | Planned | Future support |
| Gemini CLI | Planned | Future support |

## Architecture

### Directory Structure

```
pai-universal/
├── cmd/pai-universal/    # Go CLI wrapper
├── lib/                   # Core libraries
├── hooks/                 # Lifecycle hooks
├── skills/                # Skill definitions (opencode format)
├── MEMORY/
│   ├── hot/              # Session-level memory
│   ├── warm/             # Recent learnings
│   └── cold/             # Long-term memory
├── USER/
│   ├── Skills/           # User-customized skills
│   ├── Hooks/            # User hooks
│   ├── Memory/           # User memory files
│   └── Settings/         # User preferences
├── .gitlab-ci.yml        # CI/CD pipeline
└── CLAUDE.md             # This file
```

### Core Primitives

1. **Skill System** - Modular capabilities that route based on context
2. **Memory System** - 3-tier learning (hot/warm/cold) 
3. **TELOS** - Deep goal understanding files
4. **Hook System** - Lifecycle event handlers
5. **User/System Separation** - Upgrades don't overwrite user customizations

## Development

### Prerequisites

- Go 1.26+
- Git
- opencode

### Building

```bash
go build -o pai-universal ./cmd/pai-universal/
```

### Testing

```bash
go test ./...
```

### CI/CD

The project uses GitLab CI with meganerd/ci-templates:
- test: Unit tests
- build: Binary compilation
- security: Vulnerability scanning

### Issue Tracking

Issues are tracked with bd (beads):

```bash
bd list              # List all issues
bd show <id>         # Show issue details
bd create "Title"    # Create new issue
bd ready             # Show unblocked work
```

## GitLab → GitHub Mirroring

The CI pipeline includes a manual mirror job that pushes main branch to GitHub after merges. Requires `GITHUB_TOKEN` secret variable.

## License

MIT
