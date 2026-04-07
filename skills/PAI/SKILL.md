# PAI

Personal AI Infrastructure - Memory and goal management for this user.

## Overview

This skill provides access to the user's Personal AI Infrastructure, including:
- Memory system (hot/warm/cold)
- TELOS (goals, projects, mission, beliefs)
- Session insights

## When to Use

Use this skill when:
- User asks about their goals, projects, mission
- You need to reference their preferences or beliefs
- You want to log something to memory
- User mentions they're working on something (track it)
- You complete significant work (log to cold memory)
- Session context needs to be preserved for future sessions

## Key Files

- **MEMORY/hot/** - Current session working files
- **MEMORY/warm/** - Recent context (recent files, decisions)
- **MEMORY/cold/** - Long-term learnings, patterns
- **USER/TELOS/MISSION.md** - Life purpose
- **USER/TELOS/GOALS.md** - Goal tracking
- **USER/TELOS/PROJECTS.md** - Active projects
- **USER/TELOS/BELIEFS.md** - Core beliefs and preferences

## Common Tasks

### Read Memory/TELOS
```
Read USER/TELOS/GOALS.md for current goals
Read USER/TELOS/PROJECTS.md for active projects
Read MEMORY/cold/ for past learnings
```

### Log to Memory
- Write significant decisions to MEMORY/warm/
- Write learnings/patterns to MEMORY/cold/
- Update GOALS.md when new goals are discovered
- Update PROJECTS.md when new projects are mentioned

### Context for Future Sessions
- Note key files being worked on → MEMORY/hot/
- Note decisions made → MEMORY/warm/
- Extract patterns → MEMORY/cold/

## Notes

- The user prefers Go (Golang) over Node.js/npm
- They use beads (bd) for issue tracking
- They're working on pi-go (Go reimplementation of pi-mono)
- They have extensive homelab infrastructure
- See session analyzer for detailed preference/pattern extraction