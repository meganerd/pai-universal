#!/bin/bash
# PAI Hook: Session Start
# Called when entering a new opencode session in pai-universal
# Usage: ./session-start.sh ["optional prompt for complexity analysis"]

PAI_DIR="$HOME/src/Code/pai-universal"
MEMORY_DIR="$PAI_DIR/MEMORY"
SESSION_MANAGER="$PAI_DIR/cmd/session-manager"

echo "📋 PAI Session Start"
echo "==================="

# Analyze complexity if prompt provided as argument
if [ -n "$1" ]; then
    echo ""
    echo "⚡ Complexity Analysis:"
    echo "----------------------"
    
    # Get complexity score using session-manager
    cd "$PAI_DIR"
    SCORE=$(go run ./cmd/session-manager -m "$1" -score 2>&1 | grep "Complexity score:" | cut -d: -f2 | tr -d ' ')
    
    if [ -n "$SCORE" ]; then
        # Determine level based on score
        if [ "$SCORE" -le 3 ]; then
            LEVEL="Standard"
            SUGGESTION="Proceed directly - straightforward task"
        elif [ "$SCORE" -le 8 ]; then
            LEVEL="Extended"
            SUGGESTION="Consider using /prd or break into smaller tasks"
        elif [ "$SCORE" -le 16 ]; then
            LEVEL="Advanced"
            SUGGESTION="Recommend /analyze or /prd for structured approach"
        else
            LEVEL="Deep"
            SUGGESTION="Use PAI Algorithm - requires ISC breakdown"
        fi
        
        echo "  Score:     $SCORE"
        echo "  Level:     $LEVEL"
        echo "  Suggest:   $SUGGESTION"
        echo ""
        echo "  To analyze: /analyze <prompt>"
        echo "  To create PRD: /prd"
    fi
fi

# Show complexity tracking info
echo ""
echo "⚡ Complexity Tracking:"
echo "  Use: go run ./cmd/session-manager -m 'your prompt' -score"
echo "  Thresholds: Standard (0-3), Extended (4-8), Advanced (9-16), Deep (17+)"

# Check current task complexity
if [ -f "$MEMORY_DIR/hot/current-task.md" ]; then
    echo ""
    echo "Current task complexity:"
    $SESSION_MANAGER 2>/dev/null | head -5
fi

# Check for warm memory from previous sessions
if [ -d "$MEMORY_DIR/warm" ]; then
    echo ""
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
echo "📖 The Algorithm:"
echo "  For substantial tasks, say: 'Use the Algorithm'"
echo "  Read: Algorithm/v3.5.0.md"
echo ""
echo "Ready. Say what you'd like to work on."
