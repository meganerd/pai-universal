# PAI Hooks

This directory contains hooks and utilities for integrating PAI with various AI tools.

## Session Logging

Use the `log-session.sh` script to log current session work:

```bash
./tools/log-session.sh "Working on feature X"
```

This creates a session file in `MEMORY/warm/` that can be picked up by the session analyzer.

## Opencode Integration

PAI skills are automatically available to opencode when working in the pai-universal project directory. The skill is defined in `skills/PAI/SKILL.md`.

## Session Analyzer Integration

After sessions are logged to MEMORY/warm/, run the session analyzer to extract insights:

```bash
go run ./cmd/session-analyzer
```

This will:
1. Parse all session files in MEMORY/work and MEMORY/warm
2. Extract insights via LLM
3. Update memory files and TELOS in USER/TELOS/