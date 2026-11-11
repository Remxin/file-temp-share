# AGENTS.md — file-temp-share

Instructions for AI coding agents working in this repository. Read this file before changing architecture, adding packages, or implementing features.

## Product

Temporary file-sharing app (WeTransfer-like): create a **transfer**, attach **files**, share a link, download until expiry.

- Monorepo, **monolith** backend (one Go process for HTTP; separate worker process later for async jobs).
- Module name / brand: **file-temp-share**.

## Current status

| Area | Status |
|------|--------|
| Folder layout, Devbox, Taskfile, Compose (Postgres + API + frontend) | Done |
| Gin API skeleton + Swagger UI (`/swagger/index.html`) | Done |
| Initial migration (`users` / `transfers` / `files`) + user seed | Done |
| Migrations: `000001` init (incl. `blocked`) + `000002` recipients/logs; user seed | Done |
| REST CRUD, sqlc queries | Not started |
| Next.js frontend (App Router + auth/transfers UI) | In progress |
| JWT auth (single token: register/login/me/password-reset) | Done (mailer for reset not wired) |
| Local disk storage | Planned next for binaries |
| MinIO (S3) | Planned later — compose stub commented |
| RabbitMQ + worker | Setup done — topic `filetemp.emails` / queue `email.password_reset`; worker log-only (auth publish + mailer not wired) |
| Gmail mailer (personal Google account / OAuth token) | **Owner implements** — agents leave interface + env only |

## Stack

| Layer | Tech |
|-------|------|
| Backend HTTP | Go 1.23+, Gin |
| API docs | **swaggo/swag** + gin-swagger (`backend/docs`, UI at `/swagger/index.html`) |
| DB access | PostgreSQL 16, pgx/v5, **sqlc**, **golang-migrate** |
| Validation | go-playground/validator (via Gin binding / `internal/platform/validate`) |
| Logging | slog |
| Frontend | Next.js (App Router) + React + TypeScript + Tailwind |
| Task runner | Taskfile (`go-task`) |
| Dev toolbox | **Devbox** (`devbox.json`) — not a hand-written Nix flake |
| Containers | Docker Compose |
| Migrations CLI | golang-migrate `v4.18.1` built into `./bin/migrate` with `-tags postgres` (`task tools:install`) — do **not** use Nix `go-migrate` (panics on Snowflake CA) |

Pin Gin to a release compatible with the Devbox Go version (currently `gin@v1.10.1` for Go 1.23). Do not bump Gin to `@latest` without checking `go` requirements.

Do **not** introduce a second ORM (GORM, ent, etc.). Stick to sqlc + migrate.

## Repository layout

```
/
├── AGENTS.md                 # this file
├── README.md
├── .env.example
├── .envrc                    # optional direnv → devbox shellenv
├── devbox.json               # CLI toolbox
├── Taskfile.yml
├── docker-compose.yml
├── backend/
│   ├── cmd/api/              # HTTP API entrypoint
│   ├── cmd/worker/           # async consumer (RabbitMQ → log/mailer)
│   ├── internal/
│   │   ├── config/           # env / config loading
│   │   ├── server/           # Gin engine, router wiring
│   │   ├── handler/          # HTTP request layer
│   │   ├── service/          # business logic
│   │   ├── repository/       # data access wrapping sqlc
│   │   ├── dto/              # API request/response types + mappers
│   │   ├── domain/           # domain types / business errors
│   │   ├── auth/             # JWT (planned)
│   │   ├── storage/          # file blob storage (local now, MinIO later)
│   │   ├── queue/            # RabbitMQ client, publisher, consumer
│   │   ├── mailer/           # Mailer interface; Gmail impl = owner
│   │   ├── middleware/       # auth, request-id, recovery
│   │   └── platform/
│   │       ├── database/     # pool, migrate helpers
│   │       ├── httpx/        # global errors, JSON helpers
│   │       └── validate/     # validation setup
│   ├── db/
│   │   ├── migrations/       # golang-migrate *.up.sql / *.down.sql
│   │   ├── queries/          # sqlc *.sql
│   │   └── sqlc.yaml
│   ├── gen/sqlc/             # generated — do not hand-edit
│   ├── docs/                 # swag generated OpenAPI — regenerate via task api:swagger
│   ├── go.mod                # module: file-temp-share/backend
│   └── Dockerfile
├── frontend/                 # Next.js (to bootstrap later)
└── deploy/                   # future overlays (minio, rabbitmq, …)
```

## Layering rules (backend)

Strict dependency direction:

```
handler → dto
handler → service
service → domain, repository, storage, (queue)
repository → gen/sqlc (via thin wrappers)
```

- **handler**: parse/bind request, call service, map to HTTP status + DTO. No SQL.
- **service**: business rules, transactions orchestration, expiry checks, etc.
- **repository**: DB only; map sqlc models ↔ domain if needed.
- **dto**: API shapes only. **Never** expose sqlc/DB structs in JSON responses.
- **domain**: shared domain types and sentinel errors (e.g. `ErrTransferExpired`).
- **platform/httpx**: single place for error → HTTP mapping (global error handling).

### HTTP conventions

- Use correct methods: `GET` read, `POST` create, `PUT`/`PATCH` update, `DELETE` delete.
- Prefer plural resource paths, e.g. `/api/v1/transfers`, `/api/v1/transfers/:id/files`.
- Status codes: `200`/`201`/`204` success; `400` validation; `401`/`403` auth; `404` missing; `409` conflict; `410` expired transfer; `500` unexpected.
- Validation failures return structured JSON errors (field + message).

### Business rule (required beyond CRUD)

**Transfer expiry:** if `expires_at` is in the past (or status is expired), downloads and mutations that require an active transfer must fail with **410 Gone** (or a clear domain error mapped to 410). Cleanup of blobs can be async later (worker + queue).

## Domain model

Entities (5):

1. **users** — accounts (JWT later); `user_id` on transfers is required for create.
2. **transfers** — share package (`short_code`, title, `message` note, optional password, status, expires_at).
3. **files** — belongs to a transfer (`transfer_id` FK, 1—N).
4. **transfer_recipients** — emails to notify with the public link (1—N per transfer). Recipients do **not** need an account; access is public via `short_code` (+ optional transfer password).
5. **download_logs** — download history per transfer (optional `file_id` for per-file downloads).

Core relationships: `users` 1—N `transfers` 1—N `files`; `transfers` 1—N `transfer_recipients`; `transfers` 1—N `download_logs`.

## Storage

- **Now / next:** local filesystem under `UPLOAD_DIR` (Compose volume later).
- **Later:** MinIO via the same `storage.Storage` interface (put/get/delete by key).
- Keep blob I/O out of handlers; go through `internal/storage`.

## Auth

- **Single JWT** (HS256) — no refresh token pair.
- Env: `JWT_SECRET`, `JWT_TTL` (e.g. `24h`).
- Endpoints:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `GET /api/v1/auth/me` (Bearer)
  - `POST /api/v1/auth/password-reset` (issues token; **email not sent yet** — token logged for dev)
  - `POST /api/v1/auth/password-reset/confirm`
- Code: `internal/auth`, `internal/middleware`, `internal/handler` (auth), `internal/service` (auth).
- Table `password_reset_tokens` stores SHA-256 of the raw token.

## Async + mail

```
service → queue.Publisher → RabbitMQ → cmd/worker → mailer.Mailer
```

Current topology (declared idempotently by publisher/consumer):

| Element | Name |
|---------|------|
| Exchange | `filetemp.emails` (topic, durable) |
| Queue | `email.password_reset` (durable) |
| Routing key | `email.password_reset` |

- RabbitMQ runs in Compose (`task infra:up:mq` or full `infra:up`). Management UI: `http://localhost:15672` (guest/guest).
- Worker: `task worker:run` — consumes password-reset jobs and **logs** them (mailer not wired; auth does not publish yet).
- Env: `RABBITMQ_URL` (host default `amqp://guest:guest@localhost:5672/`).
- MinIO remains **commented** in `docker-compose.yml`.
- **Gmail:** send through the owner’s Google account using OAuth client credentials + refresh token (`GMAIL_*` in `.env.example`).
- Agents **must not** implement the Gmail client unless the owner asks. Provide only:
  - `mailer.Mailer` interface (e.g. `Send(ctx, msg) error`)
  - config placeholders
- Owner owns the concrete Gmail implementation under `internal/mailer`.

## Database workflow

1. Add migration pair in `backend/db/migrations/`:
   - `NNNNNN_name.up.sql` / `NNNNNN_name.down.sql`
2. Add or update queries in `backend/db/queries/*.sql`
3. Run `task db:sqlc` (generates `backend/gen/sqlc`)
4. Run `task db:migrate:up` (requires Postgres via `task infra:up`)
5. Optional: `task db:seed` loads `backend/db/seeds/*.sql` (dev only; not part of migrate)

Never edit generated sqlc output by hand. Keep seeds out of `migrations/` so sqlc schema stays clean.

## Devbox + Taskfile

- Install [Devbox](https://www.jetify.com/devbox), then `devbox shell` (or use direnv + `.envrc`).
- Tools (go, sqlc, go-task, node) come from `devbox.json`.
- **golang-migrate** is installed locally via `task tools:install` into `./bin/migrate` (postgres driver only). Devbox `init_hook` prepends `./bin` to `PATH`.
- Taskfile prefixes most tool commands with `devbox run --` via `{{.N}}`.
- Docker Compose runs **without** Devbox prefix (system Docker).

Useful tasks:

| Task | Purpose |
|------|---------|
| `task infra:up` | Build & start Postgres + RabbitMQ + API + frontend |
| `task infra:up:db` | Start only Postgres (hybrid local dev) |
| `task infra:up:mq` | Start only RabbitMQ (AMQP :5672, management :15672) |
| `task infra:down` | Stop Compose |
| `task tools:install` | Install `./bin/migrate` (postgres-only golang-migrate) |
| `task db:migrate:up` | Apply migrations |
| `task db:seed` | Load dev seed users |
| `task db:sqlc` | Generate sqlc |
| `task api:swagger` | Regenerate OpenAPI into `backend/docs` |
| `task api:run` | Run API (Swagger UI on `/swagger/index.html`) |
| `task worker:run` | Run worker (password-reset consumer; requires RabbitMQ) |
| `task devbox:check` | Verify toolbox |

### Swagger conventions

- General API metadata (`@title`, `@version`, …) lives on `backend/cmd/api/main.go`.
- Annotate handlers with swag comments; keep DTOs documented when added.
- After changing annotations, run `task api:swagger` and commit `backend/docs/*`.
- Do not hand-edit generated `docs.go` / `swagger.json` / `swagger.yaml`.

## Frontend

- Lives in `frontend/` — Next.js App Router + React + TS + Tailwind.
- Talks to backend only via REST (`frontend/src/lib/api`).
- Do not bootstrap `create-next-app` until asked; keep placeholders until then.

### Product screens — regular user (no auth detail, no admin)

Auth screens (login / register / password reset) and **admin** UI are out of scope here; admin will be specified separately. Creating transfers requires a logged-in user; **download is public**.

#### Screen map

| Screen | Route idea | Who | Purpose |
|--------|------------|-----|---------|
| Landing | `/` | public | Short product pitch + CTA to create transfer / sign in |
| Dashboard (home for logged-in) | `/app` or `/transfers` | user | List and manage own transfers |
| New transfer | `/transfers/new` | user | Wizard: files → options → confirm → share link |
| Transfer detail | `/transfers/:id` | user (owner) | Full status, files, link, actions |
| Public download | `/t/:slug` (or `/d/:token`) | public | Preview files + download |
| Expired / unavailable | same download route or dedicated | public | Clear “link expired / removed” state (HTTP 410 semantics) |
| Not found | `/404` | public | Generic 404 |
| Account settings | `/settings` | user | Minimal: change email/password (keep thin) |

Do **not** add admin screens, analytics dashboards, or folder trees unless asked.

#### Dashboard (main panel) — required content

One job: **manage the user’s transfers**. Primary CTA: **New transfer**.

**Above the list**
- Button: New transfer
- Filters: Active / Expired / All (default: Active)
- Optional: search by title

**Each transfer row/card**
- Title (or fallback name)
- Status (`active` / `expired` / …)
- File count + total size
- Expires at (`expires_at`)
- Share mode hint: link-only and/or recipient emails (when mail is enabled)
- Share link + **Copy link**
- Optional: download count
- Actions: open detail, copy link, delete (and later: extend expiry if implemented)

**Empty state**
- Short message + CTA New transfer when the user has no transfers in the current filter.

**Avoid on dashboard**
- Stats widgets, charts, “this week” strips, promo cards, unrelated settings blocks.
- Admin controls.

#### New transfer (wizard)

1. Select / upload one or more files (show names, sizes, remove).
2. Options: title, message, expiry duration, optional password, optional recipient emails (for later Gmail send).
3. Submit → create transfer → show success with share link (and copy) → link to detail.

#### Transfer detail

- Same metadata as list, plus full file list.
- Share link + copy; delete transfer; (later) resend email / extend.
- Reflect expiry: if expired, disable download-oriented actions and show expired state.

#### Public download

- Show title, message (if any), file list, total size, expiry info.
- CTA download (all or per-file — prefer per-file + “download all” if easy).
- If expired/deleted: no download buttons; explain that the link is no longer valid.

#### UX rules for agents

- Logged-in shell: simple nav (Dashboard, New transfer, Settings) — no admin entries.
- Prefer one primary action per view.
- Expiry is a first-class UI state everywhere a transfer appears.
- Wire UI to REST DTOs only; do not invent admin or auth flows in these screens beyond “user is logged in”.

## Environment

Copy `.env.example` → `.env`. Never commit secrets. Gmail tokens stay local.

## What agents should avoid

- Makefile (use Taskfile).
- Hand-written `flake.nix` (use Devbox).
- Implementing Gmail mailer, MinIO, or expanding RabbitMQ (reconnect/DLQ/auth publish) unless the task explicitly says so.
- Mixing DTO and DB models in public API responses.
- Putting business rules in handlers or SQL-only “logic” that belongs in services.
- Adding Mailhog — production mail path is Gmail via the owner’s account.
