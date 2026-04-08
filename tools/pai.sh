#!/bin/bash
# pai-opencode - Launch opencode with PAI context
#
# Usage:
#   pai-opencode                 Launch opencode in pai-universal
#   pai-opencode <project>        Launch opencode in specific project
#   pai-opencode --resume        Resume last session
#   pai-opencode --list          List sessions

OPENCODE_BIN="${HOME}/.bun/bin/opencode"
PAI_DIR="${HOME}/src/Code/pai-universal"
HOOK_DIR="$PAI_DIR/hooks"
PROJECT_DIR="${1:-$PAI_DIR}"

run_hook() {
    local hook="$1"
    if [ -f "$HOOK_DIR/$hook" ]; then
        source "$HOOK_DIR/$hook"
    fi
}

# Run session start hook (if not resuming)
if [ "$1" != "--resume" ] && [ "$1" != "-r" ]; then
    run_hook session-start.sh
    echo ""
fi

case "$1" in
    --resume|-r)
        shift
        exec $OPENCODE_BIN --resume "$@"
        ;;
    --list|-l)
        $OPENCODE_BIN session list
        ;;
    --help|-h)
        echo "pai-opencode - PAI Universal wrapper for opencode"
        echo ""
        echo "Usage:"
        echo "  pai-opencode                 Launch opencode in pai-universal"
        echo "  pai-opencode <project>       Launch opencode in specific project"
        echo "  pai-opencode --resume,-r     Resume last session"
        echo "  pai-opencode --list,-l       List sessions"
        ;;
    *)
        cd "$PROJECT_DIR"
        exec $OPENCODE_BIN .
        ;;
esac