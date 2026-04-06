#!/bin/bash
set -e

# PAI OpenCode Installation Script
# ===============================

echo "PAI OpenCode Installation"
echo "========================"
echo

# Check prerequisites
echo "Checking prerequisites..."

if ! command -v opencode &> /dev/null; then
    echo "ERROR: opencode is not installed."
    echo "Install from: https://github.com/anomalyco/opencode"
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed."
    exit 1
fi

echo "✓ opencode found"
echo "✓ Go found"
echo

# Detect installation location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PAI_DIR="${HOME}/.pai-opencode"

echo "Installation directory: ${PAI_DIR}"
echo

# Create PAI directory
mkdir -p "${PAI_DIR}"

# Copy files
echo "Installing files..."
rsync -av --exclude='.git' --exclude='.beads' \
    --exclude='go.mod' --exclude='go.sum' \
    --exclude='*.db' \
    "${SCRIPT_DIR}/" "${PAI_DIR}/"

echo "✓ Files installed to ${PAI_DIR}"
echo

# Initialize bd if not already done
if [ ! -d "${PAI_DIR}/.beads" ]; then
    echo "Initializing issue tracking..."
    cd "${PAI_DIR}" && bd init
    echo "✓ Issue tracking initialized"
fi
echo

# Create symlink for skills if needed
mkdir -p "${HOME}/.claude/skills/PAI"
if [ ! -L "${HOME}/.claude/skills/PAI" ]; then
    ln -s "${PAI_DIR}" "${HOME}/.claude/skills/PAI"
    echo "✓ Skills symlink created"
fi
echo

# Copy template TELOS files if they don't exist
for file in MISSION.md GOALS.md PROJECTS.md BELIEFS.md; do
    if [ ! -f "${PAI_DIR}/USER/TELOS/${file}" ]; then
        echo "WARNING: ${file} not found in USER/TELOS/"
    fi
done
echo

# Print completion message
echo "======================================"
echo "PAI OpenCode installed successfully!"
echo "======================================"
echo
echo "Next steps:"
echo "1. Edit USER/TELOS/MISSION.md with your life purpose"
echo "2. Edit USER/TELOS/GOALS.md with your goals"
echo "3. Edit USER/Settings/settings.json with your preferences"
echo "4. Run 'opencode' to start with PAI context"
echo
echo "For issue tracking, use:"
echo "  bd list      # List issues"
echo "  bd ready     # Show ready work"
echo
