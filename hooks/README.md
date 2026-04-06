# Hooks for PAI OpenCode

This directory contains lifecycle hooks for PAI OpenCode.

## Available Hooks

Since opencode doesn't have the same hook system as Claude Code, hooks are implemented as:

1. **Pre-session hooks** — Run before opencode starts
2. **Post-session hooks** — Run after opencode completes

## Hook Implementation

Hooks are implemented in `lib/` as shell scripts or Go utilities.

### Session Start Hook
```bash
# hooks/session-start.sh
#!/bin/bash
# Load context files, initialize memory, etc.
```

### Tool Use Hook
```bash
# hooks/tool-use.sh
#!/bin/bash
# Called after each tool use for logging/validation
```

## Example: Loading TELOS Context

```bash
#!/bin/bash
# Load TELOS files into context
TELOS_DIR="${HOME}/.pai-opencode/USER/TELOS"

if [ -d "$TELOS_DIR" ]; then
    for file in MISSION.md GOALS.md PROJECTS.md BELIEFS.md; do
        if [ -f "${TELOS_DIR}/${file}" ]; then
            echo "Loading ${file}..."
            # Add to opencode context
        fi
    done
fi
```

## Security Hooks

The security hook validates commands before execution:

```bash
#!/bin/bash
# hooks/security.sh
# Validate dangerous commands
```
