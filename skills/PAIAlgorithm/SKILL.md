# PAIAlgorithm

Personal AI Infrastructure Algorithm - Structured decision-making framework for complex tasks.

## Overview

This skill provides access to the PAI Algorithm, a structured decision-making framework that transforms complex requests into verifiable criteria and executes through structured phases.

The Algorithm helps ensure quality outcomes on non-trivial tasks by:
- Breaking down tasks into atomic, verifiable ISC (Ideal State Criteria) 
- Enforcing effort-appropriate time budgeting
- Tracking progress through 7 phases
- Announcing phase transitions via voice

## When to Use

Use this skill when:
- User requests a substantial, multi-step task
- Task complexity score is Extended (4-8) or higher
- You need to decompose a complex request into verifiable criteria
- User explicitly asks: "use the Algorithm", "/analyze", or "/prd"
- Quality must be extraordinary and time allows for structured approach

## Key Concepts

### Effort Levels

| Level | Score Range | ISC Range | When |
|-------|-------------|-----------|------|
| Standard | 0-3 | 8-16 | Normal request (DEFAULT) |
| Extended | 4-8 | 16-32 | Quality must be extraordinary |
| Advanced | 9-16 | 24-48 | Substantial multi-file work |
| Deep | 17+ | 40-80 | Complex design |

### Phases

1. **Observe** - Understand current state, gather info
2. **Think** - Analyze, decompose, generate hypotheses  
3. **Plan** - Define criteria, decide approach
4. **Build** - Execute, create, implement
5. **Verify** - Test, validate, check criteria
6. **Learn** - Extract insights, document patterns
7. **Complete** - Wrap up, summarize, hand off

## Commands

### /analyze <prompt>
Analyze task complexity and recommend effort level.

Example: `/analyze refactor the entire auth system to use JWT tokens with refresh token rotation`

Output includes:
- Complexity score
- Recommended effort level
- Suggestion for approach

### /prd or /new-prd
Create new PRD (Product Requirements Document) in WORK/ directory.

Prompts for:
- Task description (8 word imperative)
- Effort level selection
- Initial criteria (ISCs)

### /phase <name>
Transition to specified phase.

Example: `/phase verify`

Voice announces: "Entering the VERIFY phase."

### /isc <criteria>
Add ISC criterion to current PRD.

Example: `/isc JWT tokens expire within 15 minutes`

### /progress
Show current progress in active PRD.

## Key Files

- **Algorithm/v3.5.0.md** - Full Algorithm specification
- **Algorithm/PRDFORMAT.md** - PRD format reference
- **WORK/** - Active work PRDs (per-project)
- **MEMORY/warm/** - Recent session context
- **MEMORY/cold/** - Long-term learnings

## PRD Structure

PRDs are stored in `WORK/{slug}/PRD.md`:

```yaml
---
task: "8 word task description"
slug: 20260413-143000_refactor-auth-jwt
effort: extended
phase: observe
progress: 0/16
mode: interactive
started: 2026-04-13T14:30:00Z
updated: 2026-04-13T14:30:00Z
---

## Context
What was requested, why it matters, constraints

## Criteria
- [ ] ISC-1: First verifiable criterion
- [ ] ISC-2: Second verifiable criterion

## Decisions
Why certain approaches were chosen

## Verification
How to test each criterion
```

## Workflows

### Analyze Workflow
1. Accept prompt from user
2. Call session-manager with -m "<prompt>" -score
3. Display score, level, suggestion
4. If Extended+, suggest /prd

### CreatePRD Workflow
1. Ask for task description (8 words, imperative)
2. Ask for effort level
3. Generate initial PRD template
4. Guide user through ISC decomposition
5. Store in WORK/{slug}/PRD.md

### PhaseTransition Workflow
1. Accept target phase name
2. Update PRD frontmatter (phase, updated)
3. Voice announce new phase
4. Show phase-specific guidance

## Notes

- This skill is for opencode. Other harnesses have their own implementations.
- The Algorithm is a nudge, not enforcement. Users can override.
- Complexity scoring uses session-manager binary.
- Voice notifications require voice server on localhost:8888
- See Algorithm/v3.5.0.md for full specification
