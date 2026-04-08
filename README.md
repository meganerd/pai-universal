# PAI Universal - Personal AI Infrastructure

**Version:** 0.6.0  
**Status:** Full PAI with multi-harness support

## Overview

PAI Universal (Personal AI Infrastructure) brings the PAI system to work with multiple AI tools. It includes the Algorithm, session management, complexity tracking, and cross-harness sync.

## IMPORTANT: Differences from Claude Code PAI

If you're coming from Claude Code PAI, here are the key differences:

### What's the Same
- **The Algorithm** - Still available, must manually request it (say "Use the Algorithm")
- **TELOS files** - Goals, projects, mission, beliefs all work the same
- **Skills** - All PAI skill packs available (Research, Security, etc.)
- **Memory system** - 3-tier hot/warm/cold structure
- **Session analyzer** - Extracts insights from session history

### What's Different (Requires Manual Action)

| Claude Code PAI | PAI Universal (opencode) |
|-----------------|--------------------------|
| Algorithm auto-triggers on complexity | Must say "Use the Algorithm" |
| Voice notifications work | Must set up voice server separately |
| Hooks run automatically | Hooks run via pai-opencode wrapper |
| Context loads automatically | CLAUDE.md provides context |

### Workflow
1. Launch with `pai-opencode` (wrapper that runs session-start hook)
2. Tell me what you want to work on
3. For complex tasks, say "Use the Algorithm" to trigger formal process
4. Run `go run ./cmd/session-analyzer` periodically to capture learnings
5. Use `go run ./cmd/session-sync` to share context between harnesses

## Supported AI Tools

| Tool | Status | Memory Sync | Notes |
|------|--------|-------------|-------|
| opencode | ✓ Primary | ✓ | Native skill system |
| Claude Code | ✓ Active | ✓ | Reads history.jsonl + MEMORY/WORK |
| Codex | ✓ Active | ✓ | Reads logs_1.sqlite |
| Cursor | ✓ Active | ✓ | Same DB format as Codex |
| pi-go | ✓ Active | ✓ | Go rewrite of pi-mono |
| pi-mono | Reference | - | https://github.com/badlogic/pi-mono |
| Gemini CLI | Future | - | Not yet used |

## Installation

```bash
# Add to ~/.bash_aliases_local:
alias pai-opencode='bash ~/src/Code/pai-universal/tools/pai.sh'
```

Then use `pai-opencode` to launch opencode with PAI context.

## CLI Tools

### session-analyzer (Periodic)
Analyzes past sessions from all harnesses and extracts insights using LLM.

```bash
go run ./cmd/session-analyzer          # Analyze and update all
go run ./cmd/session-analyzer --dry-run -v  # Preview changes
go run ./cmd/session-analyzer -siftrank      # Auto-select best model
go run ./cmd/session-analyzer -all=false -opencode  # Specific harness
```

### session-manager (Per-Session)
Complexity detection for current task. Helps determine when to use the Algorithm.

```bash
# Score a prompt
go run ./cmd/session-manager -m "create a new web server" -score

# Show current task complexity
go run ./cmd/session-manager
```

**Complexity Thresholds:**
| Score | Level | Action |
|-------|-------|--------|
| 0-3 | Standard | Normal mode |
| 4-8 | Extended | Break into smaller tasks |
| 9-16 | Advanced | Use Algorithm (create PRD) |
| 17+ | Deep | Use Algorithm (full ISC breakdown) |

### session-sync (Cross-Harness)
Sync session context between different AI tools.

```bash
# Sync from one harness to others
go run ./cmd/session-sync --source opencode --target claude,pigo
go run ./cmd/session-sync --source claude --target opencode --dry-run
```

### notify (Notifications)
Multi-backend notification tool.

```bash
# Voice (requires voice server)
./tools/notify.sh "Task complete"

# System notification
NOTIFY_MODE=system ./tools/notify.sh "Done"
```

## Session Hooks

Hooks run automatically when using `pai-opencode`:

- **session-start.sh** - Shows goals, projects, recent sessions, complexity info
- **session-end.sh** - Logs session summary to MEMORY/warm/

## Manual Session Logging

```bash
# Log current work
./tools/log-session.sh "Working on feature X"
```

## The Algorithm

Located in `Algorithm/v3.5.0.md`. For substantial tasks:

1. Say "Use the Algorithm" or "Read Algorithm/v3.5.0.md and apply it"
2. I'll follow the 7-phase process: Observe → Reverse Engineer → Criteria → Decide → Execute → Verify → Complete
3. ISC (Ideal State Criteria) breakdown for verifiable goals

## Memory Structure

```
MEMORY/
├── hot/          # Current session files
├── warm/         # Recent sessions (parsed by session-analyzer)
└── cold/          # Long-term learnings, patterns

USER/TELOS/
├── MISSION.md     # Life purpose
├── GOALS.md       # Goal tracking
├── PROJECTS.md    # Active projects
└── BELIEFS.md     # Core beliefs and preferences
```

## License

MIT
