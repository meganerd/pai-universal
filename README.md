# PAI Universal - Personal AI Infrastructure

**Version:** 0.5.0  
**Status:** Session Analyzer + Multi-harness Support

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

## Session Analyzer

A tool to ingest session logs from various AI tools and extract learnings, preferences, and patterns to update memory, goals, and beliefs. Useful for:
- Periodically capturing learnings from your work
- Evaluating and comparing different AI tools
- Switching between harnesses frequently

### Location
`cmd/session-analyzer/`

### Usage
```bash
# Build
go build ./...

# Run (updates all harnesses by default)
go run ./cmd/session-analyzer

# Preview changes without applying
go run ./cmd/session-analyzer --dry-run -v

# Use specific model
go run ./cmd/session-analyzer -model anthropic/claude-3-5-sonnet-20241022

# With siftrank model selection (auto-selects optimal model)
go run ./cmd/session-analyzer --siftrank

# Update specific harness(es) only (disables default all)
go run ./cmd/session-analyzer -all=false -claude     # Update Claude only
go run ./cmd/session-analyzer -all=false -opencode  # Update opencode only
go run ./cmd/session-analyzer -all=false -claude -opencode  # Update multiple
```

### Environment Variables
- `PAI_BASE_DIR` - Override default base directory (defaults to `~/src/Code/pai-universal`)
- `OPENROUTER_API_KEY` - Required for LLM analysis via openrouter

### What it does
1. Parses session history from supported AI tools
2. Uses LLM to extract insights: preferences, dev patterns, projects, infrastructure, learnings, goals
3. Updates specified harnesses:
   - **pai-universal**: `MEMORY/cold/insights-*.md`, `USER/TELOS/GOALS.md`, `USER/TELOS/BELIEFS.md`
   - **Claude Code**: `~/.claude/MEMORY/WORK/insights-*.md`
   - **opencode**: `~/.local/share/opencode/storage/memory/insights-*.md`
   - **Codex/Cursor**: Same format as Claude

By default (`-all`), updates all available harnesses. Use `-all=false` with specific harness flags to target only certain ones.

### Supported Sources
- **Claude Code**: `~/.claude/history.jsonl` + `~/.claude/MEMORY/WORK/`
- **opencode**: `~/.local/share/opencode/opencode.db`
- **Codex**: `~/.codex/logs_1.sqlite` (also works for Cursor - same DB format)
- **Cursor**: `~/.cursor/logs_1.sqlite` (same schema as Codex)
- **Gemini CLI**: Not yet discovered - needs investigation

### Notes
- This is an occasional/periodic tool, not for every session
- Uses siftrank for cost-effective model selection when `--siftrank` flag is used

## GitLab → GitHub Mirroring

The CI pipeline includes a manual mirror job that pushes main branch to GitHub after merges. Requires `GITHUB_TOKEN` secret variable.

## License

MIT
