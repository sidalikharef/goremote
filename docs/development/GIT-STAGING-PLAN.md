# Git Staging Plan

## Proposed to track
- go.mod, go.sum
- All *.go source files
- web/* (embedded assets)
- docs/development/* (phase reports)
- AGENTS.md, README.md (with changes reviewed)
- .gitignore
- .github/workflows/*
- Dockerfile, docker-compose.yml, Makefile
- deploy/*
- cmd/*

## Must not track
- Binaries (goremote)
- *.auth.json, *.bootstrap
- config.json, goremote.json
- *.pem, *.key, *.crt
