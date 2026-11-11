# file-temp-share

Temporary file sharing (WeTransfer-like). Monorepo: Go/Gin API + Next.js frontend + PostgreSQL.

Architecture, layering, and feature status for AI agents: **[AGENTS.md](./AGENTS.md)**.

## Quick start

```bash
cp .env.example .env

# Full stack in Docker (Postgres :5432, API :8080, frontend :3000)
task infra:up
task db:migrate:up
task db:seed   # optional

# Or hybrid: only DB in Docker, run API/frontend via Devbox
# task infra:up:db
# task api:swagger && task api:run
# task frontend:dev
```

| Service | URL |
|---------|-----|
| Frontend | [http://localhost:3000](http://localhost:3000) |
| API / Swagger | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) |
| Postgres | `localhost:5432` |

Change host ports via `API_PORT`, `FRONTEND_PORT`, `POSTGRES_PORT` in `.env`.

## Layout

| Path | Role |
|------|------|
| `backend/` | Gin monolith (API + future worker) |
| `frontend/` | Next.js App Router |
| `docker-compose.yml` | Postgres + API + frontend; MinIO / RabbitMQ commented |
| `Taskfile.yml` | Common commands |
| `devbox.json` | Dev toolbox (go, sqlc, migrate, task, node) |
