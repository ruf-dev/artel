# Task: ArtelUI Frontend

## Goal
Build the React UI for managing Obsidian vaults, served from `pkg/client/ArtelUI`.

## Scope

### Project setup
- Initialize Bun + React + TypeScript project at `pkg/client/ArtelUI`
- Configure Vite (or Bun's bundler)
- Add `bun run dev` and `bun run build` scripts

### Features
- Vault list view
- Create vault form (name, user credentials)
- Delete vault confirmation
- Status/health display

### Backend integration
- HTTP client calling the Go backend API
- Environment-based base URL config (`.env`)

### Build integration
- `make build-ui` target: build frontend and copy `dist/` into Go backend for serving as embedded FS
- Go handler in `internal/transport/` to serve the built SPA

## Notes
- Use Bun as the package manager and runtime (not npm/yarn)
- Keep state simple (no heavy state manager needed for v1)
