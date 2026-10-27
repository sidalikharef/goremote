# Git Baseline Review

## Status
- Branch: main (commit 9c9c1ae)
- Modified: README.md
- Untracked: AGENTS.md, .gitignore, cmd/, config/, discovery/, docs/, drivers/, go.mod, go.sum, httpapi/, internal/, main.go, tv/, web/

## Proposed Tracking
All source files (go.mod, go.sum, *.go, web assets) should be tracked. Binaries, configs with secrets, .auth.json should be ignored (already in .gitignore).

## Secrets
No obvious secrets found. Auth stores are JSON files written with 0600 perms.
