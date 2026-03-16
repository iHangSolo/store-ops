# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Defaults

- Reply in **Chinese** unless I explicitly ask for English.
- No emojis.
- Do not truncate important outputs (logs, diffs, stack traces, commands,
  or critical reasoning that affects safety/correctness).

## Refactor policy (legacy code)

- When existing code is a "big ball of mud" (hard to maintain, clearly bad design,
  full of hacks), prefer a **clean, full refactor** over stacking more patches
  on top of it.
- A refactor may completely replace internal structure
  (functions, modules, classes, data flow).
- By default, try to preserve externally observable behaviour.
  If you intentionally change behaviour or protocols, you MUST:
  - Call out clearly that this is a **behaviour/protocol change**.
  - Explain why the change is necessary and which code paths/consumers are affected.
  - Update or add tests to cover the new behaviour.

## Before touching code (mandatory)

Find reuse opportunities + Trace the call/dependency chain and impact radius:

- Use semantic code search first via `codebase-retrieval` tool.
- Confirm understanding with LSP: `goToDefinition`, `findReferences`.
- Use Grep/Glob for verifying and understanding additional code snippets.

## Red lines

- No copy-paste duplication.
- Do not break existing externally observable behaviour **unless**:
  - It is part of a deliberate refactor as described in the refactor policy, and
  - You clearly document the behavioural change and its impact.
- Do not proceed with a known-wrong approach.
- Critical paths must have explicit error handling.
- Never implement "blindly": always confirm understanding via code reading + references.

## Web research (no guessing)

If something is unfamiliar or likely version-sensitive, you MUST search the web instead of guessing:

- Use Exa: `mcp__exa__web_search_exa`.

Source priority:

1. Official docs / API reference.
2. Official changelog / release notes.
3. Upstream GitHub repository docs (README, `/docs`).
4. Community posts only if necessary to fill gaps.

Version rule:

- When behaviour may differ across versions, first identify the project's version (lockfile/config),
  then search docs specifically for that version.

## Task sizing

- **Simple**
  - Criteria — single file, clear requirement, < 20 lines changed,
    clearly local impact.
  - Handling — after doing the "Before touching code" steps
    (research + impact analysis + internal three-question checklist),
    you may execute directly with minimal explanation.
  - A very short context line is enough;
    a full breakdown of the checklist is not required.

- **Medium**
  - Criteria — 2–5 files, or requires some research, or impact is not obviously local.
  - Handling — write a short plan (bullet points) → then implement.
  - Briefly surface the checklist result in the reply
    (1–3 short lines describing real issue, key reuse, and main impact).

- **Complex**
  - Criteria — architecture changes, multiple modules, high uncertainty or risk.
  - Handling — follow this workflow:
    1. **RESEARCH**: inspect code and facts only (no proposals yet).
    2. **PLAN**: present options + tradeoffs + recommendation;
       use `AskUserQuestion` actively to align with the user;
       wait for user's confirmation.
    3. **EXECUTE**: implement exactly the approved plan.
    4. **REVIEW**: self-check (tests, edge cases, cleanup).

## Git

- Do not commit unless I explicitly ask.
- Do not push unless I explicitly ask.
- Before writing a commit message, glance at a few recent commits and match the repo's style:
  - `git log -n 5 --oneline`
- If there is no obvious existing style, use this default format:
  - `<type>(<scope>): <description>`
- Before any commit: run `git diff` and confirm the exact scope of changes.
- Never force-push to `main` / `master` unless the user approves.
- Do not add attribution lines in commit messages.

## Security

- Never hardcode secrets (keys/passwords/tokens).
- Never commit `.env` files or any credentials.
- Validate user input at trust boundaries (APIs, CLIs, external data sources).

## Quality & cleanup

- Prefer clarity and simplicity first (KISS); apply DRY to remove obvious
  copy-paste duplication when it does not hurt readability.
- If you change a function signature, update **all** call sites.
- After changes:
  - Remove temporary files.
  - Remove dead/commented-out code.
  - Remove unused imports.
  - Remove debug logging that is no longer needed.
- Run the smallest meaningful verification (lint/test/build) for the parts you touched.

---

## Project Overview

This is an OpenSpec project - a spec-driven development framework for managing code changes through structured workflows. OpenSpec uses artifacts (proposal, specs, design, tasks) to document decisions before implementation.

## Prerequisites

- **OpenSpec CLI** (`openspec-cn`) must be installed

## Core Commands

All OpenSpec commands are invoked via skills:

| Command | Purpose |
|---------|---------|
| `/opsx:propose` | Create change and generate all artifacts in one step |
| `/opsx:new` | Start new change, step through artifacts one-by-one |
| `/opsx:continue` | Continue working on existing change |
| `/opsx:apply` | Implement tasks from a change |
| `/opsx:verify` | Verify implementation matches artifacts |
| `/opsx:archive` | Archive completed change |
| `/opsx:explore` | Thinking/exploration mode (no code changes) |
| `/opsx:ff` | Fast-forward: create all artifacts at once |
| `/opsx:onboard` | Guided tutorial for first-time users |

## Architecture

```
openspec/
├── config.yaml        # Project config (schema: spec-driven)
├── specs/             # Main specifications (persistent)
├── changes/           # Active changes
│   ├── <change-name>/ # Each change is a directory
│   │   ├── .openspec.yaml
│   │   ├── proposal.md
│   │   ├── design.md
│   │   ├── tasks.md
│   │   └── specs/<capability>/spec.md
│   └── archive/       # Completed changes (YYYY-MM-DD-<name>)
```

## Workflow

1. **Explore** (`/opsx:explore`) - Think through the problem
2. **Create** (`/opsx:new` or `/opsx:propose`) - Create change with artifacts
3. **Implement** (`/opsx:apply`) - Work through tasks
4. **Verify** (`/opsx:verify`) - Check implementation matches specs
5. **Archive** (`/opsx:archive`) - Move to archive with date prefix

## Key CLI Commands

```bash
openspec-cn list --json                    # List all changes
openspec-cn status --change "<name>" --json # Check change status
openspec-cn new change "<name>"             # Create new change
openspec-cn instructions <artifact> --change "<name>" --json  # Get artifact template
```

## Artifact Types (spec-driven schema)

- **proposal.md** - Why we're doing this change
- **specs/<capability>/spec.md** - Detailed requirements (When/Then format)
- **design.md** - Technical decisions and approach
- **tasks.md** - Implementation checklist with `- [ ]` items

## Skills Location

Skills are defined in `.claude/skills/openspec-*/SKILL.md` and commands in `.claude/commands/opsx/*.md`.