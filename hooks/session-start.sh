#!/bin/bash
# PAI Hook: Session Start
# Called when entering a new opencode session in pai-universal

PAI_DIR="$HOME/src/Code/pai-universal"
MEMORY_DIR="$PAI_DIR/MEMORY"

echo "📋 PAI Session Start"
echo "==================="

# Check for warm memory from previous sessions
if [ -d "$MEMORY_DIR/warm" ]; then
    echo "Recent sessions in MEMORY/warm/:"
    ls -t "$MEMORY_DIR/warm/" | head -5 | while read f; do
        echo "  - $f"
    done
fi

# Check for cold patterns
if [ -d "$MEMORY_DIR/cold" ]; then
    echo ""
    echo "Known patterns in MEMORY/cold/:"
    ls -t "$MEMORY_DIR/cold/" | head -3 | while read f; do
        echo "  - $f"
    done
fi

# Load TELOS context
echo ""
echo "📌 Active Goals:"
if [ -f "$PAI_DIR/USER/TELOS/GOALS.md" ]; then
    grep "^##" "$PAI_DIR/USER/TELOS/GOALS.md" 2>/dev/null | head -5
fi

echo ""
echo "📦 Active Projects:"
if [ -f "$PAI_DIR/USER/TELOS/PROJECTS.md" ]; then
    grep "^###" "$PAI_DIR/USER/TELOS/PROJECTS.md" 2>/dev/null | head -5
fi

echo ""
echo "Ready. Say what you'd like to work on."
