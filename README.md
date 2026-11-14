# file-temp-share

Temporary file sharing (WeTransfer-like). Monorepo: Go/Gin API + Next.js frontend + PostgreSQL.

Architecture, layering, and feature status for AI agents: **[AGENTS.md](./AGENTS.md)**.

## Quick start

```bash
cp .env.example .env

# Full stack in Docker (Postgres, RabbitMQ, API, worker, frontend)
task infra:up
task db:migrate:up
task db:seed   # optional

# Or hybrid: infra in Docker, processes via Devbox
# task infra:up:db && task infra:up:mq
# task api:swagger && task api:run
# task worker:run
# task frontend:dev
```

| Service | URL |
|---------|-----|
| Frontend | [http://localhost:3000](http://localhost:3000) |
| API / Swagger | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) |
| Postgres | `localhost:5432` |
| MinIO S3 API | `localhost:9000` |
| MinIO console | [http://localhost:9001](http://localhost:9001) (default `minioadmin` / `minioadmin`) |
| RabbitMQ AMQP | `localhost:5672` |
| RabbitMQ management | [http://localhost:15672](http://localhost:15672) (guest/guest) |

Change host ports via `API_PORT`, `FRONTEND_PORT`, `POSTGRES_PORT`, `MINIO_API_PORT`, `MINIO_CONSOLE_PORT`, `RABBITMQ_PORT`, `RABBITMQ_MANAGEMENT_PORT` in `.env`.

Worker (Compose service `worker`, image `backend/Dockerfile.worker`) sends password-reset mail via Gmail SMTP. Put `GMAIL_APP_LOGIN` / `GMAIL_APP_PASSWORD` in project `.env` so Compose can inject them.

## Layout

| Path | Role |
|------|------|
| `backend/` | Gin monolith (API + worker) |
| `frontend/` | Next.js App Router |
| `docker-compose.yml` | Postgres + MinIO + RabbitMQ + API + worker + frontend |
| `Taskfile.yml` | Common commands |
| `devbox.json` | Dev toolbox (go, sqlc, migrate, task, node) |
