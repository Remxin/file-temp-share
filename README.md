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

# Generate Swagger docs, then run API
task api:swagger
task api:run
```

Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Layout

| Path | Role |
|------|------|
| `backend/` | Gin monolith (API + future worker) |
| `frontend/` | Next.js (placeholder until bootstrapped) |
| `docker-compose.yml` | Postgres now; MinIO / RabbitMQ commented |
| `Taskfile.yml` | Common commands |
| `devbox.json` | Dev toolbox (go, sqlc, migrate, task, node) |
