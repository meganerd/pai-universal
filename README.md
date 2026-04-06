# PAI Universal - Personal AI Infrastructure

**Version:** 0.5.0  
**Status:** Session Analyzer + Multi-harness Support

## Overview

PAI Universal (Personal AI Infrastructure) is a goal-oriented AI assistant framework that works across multiple AI tools. Unlike stateless AI assistants, PAI knows your goals, remembers your preferences, and continuously improves its understanding of how you work.

### Supported AI Tools

| Tool | Status | Memory Sync | Notes |
|------|--------|-------------|-------|
| opencode | ✓ Primary | ✓ | Native skill system |
| Claude Code | ✓ Active | ✓ | Reads history.jsonl + MEMORY/WORK |
| Codex | ✓ Active | ✓ | Reads logs_1.sqlite |
| Cursor | ✓ Active | ✓ | Same DB format as Codex |
| pi-go | ✓ Active | ✓ | Go rewrite of pi-mono |
| pi-mono | Reference | - | https://github.com/badlogic/pi-mono |
| Gemini CLI | Future | - | Not yet used |

## Session Analyzer

A tool to ingest session logs from various AI tools and extract learnings, preferences, and patterns to update memory, goals, and beliefs.

### Why use it?
- **Capture learnings** - Periodically extract what you've done across all AI tools
- **Compare tools** - See what you do in each harness
- **Multi-harness workflow** - Keep context when switching between tools

### Quick Start
```bash
# Build all tools
go build ./...

# Run session analyzer (updates all harnesses by default)
go run ./cmd/session-analyzer

# Preview what would change
go run ./cmd/session-analyzer --dry-run -v

# Update specific harnesses only
go run ./cmd/session-analyzer -all=false -opencode
```

### Output Targets

**pai-universal (local):**
- `MEMORY/cold/insights-YYYYMMDD.md` - New learnings
- `USER/TELOS/GOALS.md` - Inferred goals
- `USER/TELOS/BELIEFS.md` - Technology preferences

**Harnesses:**
- Claude Code: `~/.claude/MEMORY/WORK/insights-*.md`
- opencode: `~/.local/share/opencode/storage/memory/insights-*.md`
- Codex/Cursor: Same format as Claude
- pi-go: `~/.local/share/pi-go/memory/insights-*.md`

### Configuration
- `PAI_BASE_DIR` - Override default (default: `~/src/Code/pai-universal`)
- `OPENROUTER_API_KEY` - Required for LLM analysis
- `-model` - Set LLM model (default: `google/gemini-2.0-flash-001`)
- `-siftrank` - Use siftrank for auto model selection

### Notes
- This is an occasional/periodic tool, not for every session
- Source data is read-only; creates new insight files
- Handles missing sources gracefully

## License

MIT
