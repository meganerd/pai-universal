#!/bin/bash
# PAI Hook: Session End
# Called when ending an opencode session in pai-universal

PAI_DIR="$HOME/src/Code/pai-universal"
WORK_DIR="$PAI_DIR/MEMORY/warm"
SESSION_FILE="$WORK_DIR/session-$(date +%Y%m%d-%H%M%S).md"

echo "📋 PAI Session End"
echo "=================="

# Get session summary
SUMMARY="${1:-$(git -C "$PAI_DIR" log --oneline -3 2>/dev/null | head -3)}"

cat > "$SESSION_FILE" << EOF
# Session: $(date +%Y-%m-%d %H:%M)

## Summary
$SUMMARY

## Files Changed
$(git -C "$PAI_DIR" diff --stat HEAD 2>/dev/null | tail -10)

---
*Logged via PAI session-end hook*
EOF

echo "Session logged to: $SESSION_FILE"
echo ""
echo "Run session analyzer to extract insights:"
echo "  go run ./cmd/session-analyzer"
