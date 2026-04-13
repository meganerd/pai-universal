# PAI OpenCode Skill System

This directory contains skill definitions that integrate with opencode's skill system.

## Available Skills

The following PAI skills are already available in your opencode setup at `~/.claude/skills/`:

- **Research** — Multi-agent research with quick/standard/extensive/deep modes
- **Security** — Recon, web assessment, prompt injection testing, security news
- **Media** — Image generation, diagrams, Remotion video
- **Investigation** — OSINT and people search
- **Thinking** — First principles, council debates, red team, brainstorming, science
- **Telos** — Life OS for goals, beliefs, wisdom
- **USMetrics** — 68 US economic indicators
- **Scraping** — Web scraping with Bright Data and Apify
- **ContentAnalysis** — Wisdom extraction from videos, podcasts, articles
- **Agents** — Custom agent composition
- **Utilities** — CLI generation, skill scaffolding, Fabric patterns

### PAIAlgorithm (pai-universal only)

Located in `skills/PAIAlgorithm/` — available when working in pai-universal:

- **PAIAlgorithm** — Structured decision-making framework for complex tasks
  - `/analyze <prompt>` — Score complexity, recommend effort level
  - `/prd` — Create new PRD in WORK/ directory
  - `/phase <name>` — Transition Algorithm phase
  - `/isc <criteria>` — Add ISC criterion to PRD

See `PAIAlgorithm/SKILL.md` for full documentation.

## Customization

To customize any skill, create overrides in `USER/Skills/`:

```
USER/Skills/Research/
├── SKILL.md          # Override main skill behavior
├── Workflows/        # Add custom workflows
└── Templates/        # Add custom templates
```

## Adding New Skills

1. Create directory: `skills/<SkillName>/`
2. Add `SKILL.md` with YAML frontmatter
3. Optionally add `Workflows/` and `Tools/` subdirectories

See `~/.claude/skills/Utilities/CreateSkill/` for scaffolding help.
