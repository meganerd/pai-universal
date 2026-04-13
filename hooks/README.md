# PAI Hooks

This directory contains hooks and utilities for integrating PAI with various AI tools.

## Hooks

### Session Start Hook
`session-start.sh` - Runs when entering a PAI session

Usage:
```bash
./session-start.sh                    # Basic info only
./session-start.sh "your prompt"     # Analyze complexity
```

Shows:
- Recent sessions from MEMORY/warm/
- Known patterns from MEMORY/cold/
- Active goals from USER/TELOS/GOALS.md
- Active projects from USER/TELOS/PROJECTS.md
- Complexity analysis if prompt provided

### Session End Hook  
`session-end.sh` - Runs when ending a PAI session

Logs:
- Session summary
- Files changed (git diff)
- Timestamp

Output goes to MEMORY/warm/ for pickup by session analyzer.

## Notification System

### notify.sh
Multi-backend notification tool:

```bash
# Voice (requires voice server on localhost:8888)
./tools/notify.sh "Entering the Algorithm"

# System notifications
NOTIFY_MODE=system ./tools/notify.sh "Task complete"

# Log only (fallback)
NOTIFY_MODE=log ./tools/notify.sh "Checkpoint reached"
```

Environment variables:
- `VOICE_SERVER` - Voice server URL (default: localhost:8888)
- `NOTIFY_MODE` - voice, system, or log

## Session Logging

Use the log-session script:
```bash
./tools/log-session.sh "Working on feature X"
```

This creates a session file in `MEMORY/warm/`.

## Opencode Integration

PAI skills are automatically available to opencode when working in the pai-universal project directory.

## Session Analyzer

After sessions are logged, run:
```bash
go run ./cmd/session-analyzer
```

This parses MEMORY/warm/ and extracts insights to MEMORY/cold/ and TELOS files.