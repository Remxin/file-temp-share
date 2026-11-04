# file-temp-share

Temporary file sharing (WeTransfer-like). Monorepo: Go/Gin API + Next.js frontend + PostgreSQL.

Architecture, layering, and feature status for AI agents: **[AGENTS.md](./AGENTS.md)**.

## Quick start

```bash
cp .env.example .env

# Toolchain (optional but recommended)
# https://www.jetify.com/devbox
devbox shell

# Postgres
task infra:up

# API stub
task api:run
```

## Layout

| Path | Role |
|------|------|
| `backend/` | Gin monolith (API + future worker) |
| `frontend/` | Next.js (placeholder until bootstrapped) |
| `docker-compose.yml` | Postgres now; MinIO / RabbitMQ commented |
| `Taskfile.yml` | Common commands |
| `devbox.json` | Dev toolbox (go, sqlc, migrate, task, node) |

## Requirements (course)

REST API, PostgreSQL, ≥3 entities with a relation, CRUD, layered backend, DTO, validation, global errors, ≥1 business rule beyond CRUD, correct HTTP codes — see AGENTS.md for how this repo maps to those requirements.
