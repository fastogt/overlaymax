# CLAUDE.md

## Project Overview

**overlaymax** — Sports overlay data service managing football/basketball event metadata and graphical overlays via WebSocket real-time updates. Uses Pogreb embedded KV store.

## Build & Development Commands

```bash
go build .                              # Build binary
golangci-lint run ./...                 # Lint (run from root)
```

## Architecture

Application (root-level module, no `src/` directory). Module: `backend`. Entry point: `main.go`

Default port: **8180**. Config: `configs/overlaymax.conf`

**Note:** Uses `gorilla/mux` for routing (not chi). Root-level Go module — no `src/` subdirectory.

### Key Packages

- `app/` — core application: routes, config, server
- `app/models/` — data models (football, basketball, overlay)
- `app/store/` — Pogreb database abstraction
- `app/updates/` — WebSocket update manager

### Key Dependencies

- `github.com/gorilla/mux` — HTTP routing
- `github.com/gorilla/websocket` — WebSocket
- `github.com/akrylysov/pogreb` — embedded KV database
- `gitlab.com/fastogt/gofastogt` — core framework
