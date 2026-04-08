#!/bin/bash
# PAI Session Logger - Log current session work to memory
# Usage: ./log-session.sh "What did you work on?"

set -e

MEMORY_DIR="${PAI_BASE_DIR:-$HOME/src/Code/pai-universal}/MEMORY"
WORK_DIR="$MEMORY_DIR/warm"
SESSION_FILE="session-$(date +%Y%m%d-%H%M%S).md"

# Get session description from args or prompt
if [ -n "$1" ]; then
    DESC="$1"
else
    read -p "What did you work on? " DESC
fi

cat > "$WORK_DIR/$SESSION_FILE" << EOF
# Session: $(date +%Y-%m-%d)

## Work Done

$DESC

## Files Changed

$(git status --short 2>/dev/null || echo "(not a git repo)")

## Next Steps

_TBD_

---
*Logged from opencode session*
EOF

echo "Logged to: $WORK_DIR/$SESSION_FILE"