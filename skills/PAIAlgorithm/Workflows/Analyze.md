# Analyze Workflow

Analyzes task complexity and recommends effort level.

## Trigger

`/analyze <prompt>` or skill invocation for complexity assessment

## Steps

### 1. Accept Prompt
Get the prompt to analyze. Can be passed as argument or extracted from conversation.

### 2. Run Complexity Scoring
Call the session-manager binary:

```bash
go run ./cmd/session-manager -m "<prompt>" -score
```

Or if in skill context, use the tool call to run the binary.

### 3. Display Results
Show:
- **Score**: Integer 0+
- **Level**: Standard | Extended | Advanced | Deep
- **Suggestion**: What to do next

### 4. Recommend Next Steps

| Level | Suggestion |
|-------|-------------|
| Standard (0-3) | Proceed directly - task is straightforward |
| Extended (4-8) | Consider breaking into smaller tasks or use /prd |
| Advanced (9-16) | Recommend using Algorithm - create PRD |
| Deep (17+) | Use Algorithm - requires ISC breakdown |

### 5. Offer Actions

If Extended+:
- Offer to create PRD: `/prd`
- Offer to transition to Observe phase: `/phase observe`
- Allow override: "proceed anyway"

## Output Format

```
═══════════════════════════════════════
  COMPLEXITY ANALYSIS
═══════════════════════════════════════
  Score:    6
  Level:    Extended
  
  Suggestion: Consider breaking into smaller tasks or use /prd
  
  To proceed:
    /prd        → Create PRD with ISC criteria
    /phase observe → Start Algorithm from Observe phase
    "proceed anyway" → Skip Algorithm, work directly
═══════════════════════════════════════
```

## Notes

- The complexity algorithm counts: file mentions, directory patterns, keywords (rewrite, migrate, refactor, design), testing/CI mentions, prompt length
- Higher scores don't mean "can't do" - they mean "should structure approach"
- User always has override option
