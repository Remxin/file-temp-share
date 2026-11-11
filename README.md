# file-temp-share

Temporary file sharing (WeTransfer-like). Monorepo: Go/Gin API + Next.js frontend + PostgreSQL.

Architecture, layering, and feature status for AI agents: **[AGENTS.md](./AGENTS.md)**.

## Quick start

```bash
cp .env.example .env

# Full stack in Docker (Postgres :5432, RabbitMQ :5672, API :8080, frontend :3000)
task infra:up
task db:migrate:up
task db:seed   # optional

# Or hybrid: only DB (+ optionally MQ) in Docker, run API/frontend/worker via Devbox
# task infra:up:db
# task infra:up:mq
# task api:swagger && task api:run
# task worker:run
# task frontend:dev
```

| Service | URL |
|---------|-----|
| Frontend | [http://localhost:3000](http://localhost:3000) |
| API / Swagger | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) |
| Postgres | `localhost:5432` |
| RabbitMQ AMQP | `localhost:5672` |
| RabbitMQ management | [http://localhost:15672](http://localhost:15672) (guest/guest) |

Change host ports via `API_PORT`, `FRONTEND_PORT`, `POSTGRES_PORT`, `RABBITMQ_PORT`, `RABBITMQ_MANAGEMENT_PORT` in `.env`.

Worker (`task worker:run`) consumes `email.password_reset` and sends mail via Gmail SMTP (`GMAIL_APP_LOGIN` / `GMAIL_APP_PASSWORD` in `.env` or `.env.local`).

## Layout

| Path | Role |
|------|------|
| `backend/` | Gin monolith (API + worker) |
| `frontend/` | Next.js App Router |
| `docker-compose.yml` | Postgres + RabbitMQ + API + frontend; MinIO commented |
| `Taskfile.yml` | Common commands |
| `devbox.json` | Dev toolbox (go, sqlc, migrate, task, node) |
