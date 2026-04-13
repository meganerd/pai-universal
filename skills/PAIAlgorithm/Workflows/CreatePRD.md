# CreatePRD Workflow

Creates a new PRD (Product Requirements Document) for Algorithm-based work.

## Trigger

`/prd` or `/new-prd` command

## Prerequisites

- Task complexity is Extended+ (score 4+), OR
- User explicitly requests structured approach

## Steps

### 1. Collect Task Information

Ask user for:
1. **Task description** (8 words, imperative mood)
   - Example: "Add JWT refresh token rotation"
   - Should describe deliverable, not process

2. **Effort level** (default: Extended)
   - Standard (0-3): Simple tasks
   - Extended (4-8): Quality matters
   - Advanced (9-16): Multi-file work
   - Deep (17+): Complex design

3. **Initial criteria** (at least 3)
   - Ask user for their success criteria
   - Help decompose into atomic ISCs

### 2. Generate Slug

Create unique identifier:
```
YYYYMMDD-HHMMSS_kebab-task-description
```

### 3. Create WORK Directory

```
WORK/
└── {slug}/
    └── PRD.md
```

### 4. Write PRD Template

```yaml
---
task: "{8 word task}"
slug: {YYYYMMDD-HHMMSS_kebab-task}
effort: {level}
phase: observe
progress: 0/{estimated_isc_count}
mode: interactive
started: {ISO timestamp}
updated: {ISO timestamp}
---

## Context

What was requested:
- 

Why it matters:
- 

Key constraints:
- 

## Criteria

- [ ] ISC-1: 
- [ ] ISC-2: 
- [ ] ISC-3: 

## Decisions

### Approach 1
- Pro: 
- Con: 

### Approach 2  
- Pro:
- Con:

Chosen approach:

## Verification

For each ISC, document how to verify:
- 

## Notes

Additional context and notes.
```

### 5. Guide ISC Decomposition

Help user break down their criteria into atomic ISCs:

**The Splitting Test:**
- If criteria contains "and" → split into separate criteria
- Can part A pass while part B fails? → separate criteria
- "All" or "every" → enumerate what "all" means

**Decomposition by domain:**
| Domain | Decompose per... |
|--------|-----------------|
| UI/Visual | Element, state, breakpoint |
| Data/API | Field, validation rule, edge case |
| Logic/Flow | Branch, transition, boundary |
| Content | Section, format, tone |
| Infrastructure | Service, config, permission |

### 6. Update Progress

After criteria finalized:
- Count total ISCs
- Update progress in frontmatter: `progress: 0/{count}`

### 7. Announce Phase

Voice announce: "PRD created. Entering the OBSERVE phase."

## Post-Creation

After PRD is created:
- Transition to Observe phase
- Begin gathering context
- User can edit PRD directly or ask for help

## Notes

- PRD is the single source of truth
- Only the AI writes to PRD (via Edit/Write tools)
- Hooks read PRD to sync state (but don't write)
- User can abandon PRD anytime with "cancel" or "skip"
