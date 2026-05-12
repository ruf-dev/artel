# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**artel** is a Go service that automates Obsidian in-cloud vault creation. It provisions CouchDB instances, creates
users and resources within them, and exposes a React UI at `pkg/client/ArtelUI` (Bun + React). Generated
by [RedSock CLI (rscli)](https://github.com/Red-Sock/rscli).

### Frontend (pkg/client/ArtelUI)

### Code Generation

```bash
moti i          # Install protoc dependencies if new added to moti.yaml
moti g          # Generate protos (preferred — runs protoc via moti.yaml config)
sqlc generate   # Generate Go from SQL queries (internal/repository/pg/queries → internal/repository/pg/generated)
make lint       # golangci-lint run ./...
```

## Architecture

If touching Go backend files, read [docs/architecture.md](docs/architecture.md) for rscli structure, configuration,
planned layers, service layout, and key dependencies.

## Go Coding Rules

If editing any `.go` files, read and follow [docs/go-style.md](docs/go-style.md).