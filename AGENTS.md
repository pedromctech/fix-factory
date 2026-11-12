# AGENTS.md

Context for the factory's agents: replace this with how THIS project works (Claude Code can write it).

## Commands
- Test: the TEST_CMD in factory.conf (it must pass before any PR)

## Conventions
- (how the code is organised, naming, error handling, what a good test looks like)

## Factory rules
- Phase state lives in `.factory/` as markdown files. Never put plan or review text in code comments.
- Only the paths in SCOPE (factory.conf) may change. Do not touch factory.sh, doctor.sh, lib/, gates/,
  .claude/, dashboard/, factory.conf or this file.
