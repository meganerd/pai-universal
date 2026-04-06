# PAI Universal - CLAUDE.md

**Version:** 0.1.0-alpha  
**Purpose:** Universal Personal AI Infrastructure supporting multiple AI tools

## Project Overview

This is PAI (Personal AI Infrastructure) adapted to work with multiple AI tools, starting with opencode.

## Supported AI Tools

| Tool | Status | Notes |
|------|--------|-------|
| opencode | ✓ Primary | Native skill system |
| Claude Code | Legacy | Original PAI |
| Cursor | Planned | Future support |
| Codex | Planned | Future support |
| Gemini CLI | Planned | Future support |

## Key Differences from Claude Code PAI

| Component | Claude Code PAI | PAI OpenCode |
|-----------|-----------------|--------------|
| Skills | Claude Code hooks | opencode Skill tool |
| Memory | 3-tier (hot/warm/cold) | Local filesystem |
| Hooks | Claude Code hooks | opencode lifecycle |
| Context | CLAUDE.md + hooks | Context files + skills |
| Installation | install.sh | Manual setup |

## Skills Available

All PAI Packs are already available as opencode skills:
- `~/.claude/skills/Research/`
- `~/.claude/skills/Security/`
- `~/.claude/skills/Media/`
- `~/.claude/skills/Telos/`
- `~/.claude/skills/Thinking/`
- etc.

## Working with This Project

### Issue Tracking
Use `bd` for issue tracking:
```bash
bd list           # Show all issues
bd ready          # Show unblocked work
bd show <id>      # Show issue details
```

### Development
```bash
go build ./...    # Build Go utilities
go test ./...     # Run tests
```

### Context Files
- `USER/TELOS/MISSION.md` — Life purpose
- `USER/TELOS/GOALS.md` — Goals tracking
- `USER/TELOS/PROJECTS.md` — Active projects
- `USER/Settings/settings.json` — User preferences

## Phase 2: Memory Migration

After Phase 1 is complete, we will import:
- Lessons learned from `~/.claude/PAI/MEMORY/`
- TELOS content from `~/.claude/PAI/USER/TELOS/`
- Custom skill overrides from `~/.claude/PAI/USER/SKILLCUSTOMIZATIONS/`
