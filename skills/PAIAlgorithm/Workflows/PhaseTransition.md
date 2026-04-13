# PhaseTransition Workflow

Transitions between Algorithm phases with voice announcement and PRD update.

## Trigger

`/phase <name>` command

## Valid Phases

1. **observe** - Understand current state, gather info
2. **think** - Analyze, decompose, generate hypotheses
3. **plan** - Define criteria, decide approach
4. **build** - Execute, create, implement
5. **verify** - Test, validate, check criteria
6. **learn** - Extract insights, document patterns
7. **complete** - Wrap up, summarize, hand off

## Steps

### 1. Validate Phase Name

Accept phase name (case-insensitive). Map to canonical:
- obs, o → observe
- think, t → think
- plan, p → plan
- build, b → build
- verify, v → verify
- learn, l → learn
- complete, c, done → complete

### 2. Check for Active PRD

If no active PRD exists:
- Warn user: "No active PRD. Consider /prd first."
- Allow transition anyway if user insists

### 3. Update PRD Frontmatter

Edit PRD.md frontmatter:
```yaml
phase: {new_phase}
updated: {ISO timestamp}
```

### 4. Voice Announcement

Send voice notification:
```bash
curl -s -X POST http://localhost:8888/notify \
  -H "Content-Type: application/json" \
  -d '{"message": "Entering the {PHASE} phase.", "voice_id": "a1TnjruAs5jTzdrjL8Vd", "voice_enabled": true}'
```

### 5. Show Phase Guidance

Display phase-specific tips:

| Phase | Guidance |
|-------|----------|
| **observe** | Gather context: read files, check docs, understand current state |
| **think** | Analyze options, decompose problem, generate approaches |
| **plan** | Define ISC criteria, choose approach, document decisions |
| **build** | Implement, write code, create tests |
| **verify** | Test each criterion, validate against ISCs |
| **learn** | Document what worked, what didn't, patterns to reuse |
| **complete** | Summarize work, hand off to user |

### 6. Update Progress (if needed)

If transitioning to verify or complete:
- Show progress: "X/Y criteria complete"
- List incomplete criteria

## Output Format

```
═══════════════════════════════════════
  PHASE TRANSITION
═══════════════════════════════════════
  Previous:  observe
  Current:   think
  
  ═════════════════════════════════════
  Think Phase Guidance:
  ═════════════════════════════════════
  - Analyze the problem from multiple angles
  - Consider first principles
  - Generate hypotheses about root causes
  - Decompose into smaller sub-problems
  
  Use /isc to add criteria as you discover them.
═══════════════════════════════════════
```

## Backward Transitions

Backward transitions are allowed:
- verify → build (fix issues found)
- think → observe (need more context)
- complete → verify (criteria failed)

Forward transitions:
- Must complete current phase before moving forward
- Can skip phases with user permission

## Notes

- Only the primary agent makes voice announcements
- Background agents skip voice calls
- Phase name is stored in PRD frontmatter
- Progress counter shows checked/total ISCs
