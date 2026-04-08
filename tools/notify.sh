#!/bin/bash
# PAI Notification Tool
# Usage: ./notify.sh "message" [voice_id]
#
# Supports multiple backends:
# - local voice server (localhost:8888)
# - system notifications (notify-send)
# - log only

VOICE_SERVER="${VOICE_SERVER:-localhost:8888}"
NOTIFY_MODE="${NOTIFY_MODE:-voice}"  # voice, system, log

case "$NOTIFY_MODE" in
    voice)
        if curl -s --connect-timeout 1 "http://$VOICE_SERVER" >/dev/null 2>&1; then
            VOICE_ID="${2:-a1TnjruAs5jTzdrjL8Vd}"
            curl -s -X POST "http://$VOICE_SERVER/notify" \
                -H "Content-Type: application/json" \
                -d "{\"message\": \"$1\", \"voice_id\": \"$VOICE_ID\", \"voice_enabled\": true}" \
                >/dev/null 2>&1
        fi
        ;;
    system)
        if command -v notify-send >/dev/null 2>&1; then
            notify-send "PAI" "$1"
        fi
        ;;
    log)
        echo "[PAI NOTIFY] $(date): $1"
        ;;
esac
