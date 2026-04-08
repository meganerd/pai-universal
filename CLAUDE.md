# PAI Universal - CLAUDE.md

**Version:** 0.6.0
**Purpose:** Universal Personal AI Infrastructure supporting multiple AI tools

## Project Overview

This is PAI (Personal AI Infrastructure) adapted to work with multiple AI tools, starting with opencode.

## The Algorithm

**Read `Algorithm/v3.5.0.md` before any substantial work.**

The Algorithm is the core decision-making framework. It defines:
- Effort levels (Standard, Extended, Advanced, Deep, Comprehensive)
- 7 phases: Observe, Reverse Engineer, Criteria, Decide, Execute, Verify, Complete
- ISC (Ideal State Criteria) decomposition methodology
- Voice notifications for phase transitions

**Usage:** For any non-trivial task, apply the Algorithm. It transforms requests into verifiable criteria and executes through structured phases.

## Supported AI Tools

| Tool | Status | Notes |
|------|--------|-------|
| opencode | ✓ Primary | Native skill system |
| Claude Code | ✓ Active | Reads history.jsonl + MEMORY/WORK |
| Codex | ✓ Active | Reads logs_1.sqlite |
| Cursor | ✓ Active | Same DB format as Codex |
| pi-go | ✓ Active | Go rewrite of pi-mono |

## CLI Tools (for manual use)

- `go run ./cmd/session-analyzer` - Extract insights from session history
- `go run ./cmd/session-manager -m "prompt" -score` - Check task complexity
- `go run ./cmd/session-sync --source X --target Y` - Cross-harness sync

## Running PAI Universal

Use the `pai-opencode` alias to launch with PAI context:
```bash
pai-opencode              # Launch in pai-universal
pai-opencode ~/src/Code/pi-go  # Launch in specific project
pai-opencode --resume      # Resume last session
```

The session start hook will show:
- Current task complexity
- Recent sessions in MEMORY/warm/
- Active goals and projects from TELOS
- Reminder to use the Algorithm for complex tasks

## Context Files
- `USER/TELOS/MISSION.md` — Life purpose
- `USER/TELOS/GOALS.md` — Goals tracking
- `USER/TELOS/PROJECTS.md` — Active projects
- `USER/Settings/settings.json` — User preferences
