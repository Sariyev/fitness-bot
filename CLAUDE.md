# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Telegram Mini App fitness coaching platform. Three Go binaries (`bot`, `admin`, `webapp`) share a single Postgres DB and the same `internal/` packages. The Vue 3 frontend in [web/](web/) is built to static files and served by the `webapp` binary. Payments go through Robokassa, with a Dummy provider as a fallback for local dev.

## Common commands

- **Start Postgres only** (host port **5444**, container port 5432): `docker-compose up -d db`
- **Run a service from source**: `go run ./cmd/bot` | `go run ./cmd/admin` | `go run ./cmd/webapp`
- **Frontend dev server**: `cd web && npm install && npm run dev`
- **Frontend production build** (TS check + bundle to [web/dist/](web/dist/)): `cd web && npm run build` — runs `vue-tsc --noEmit && vite build`
- **Build & run all services in containers**: `docker-compose up -d --build`
- **Migrations**: run automatically on every service startup via `database.RunMigrations` (see [cmd/webapp/main.go](cmd/webapp/main.go)). To add one, drop a numbered `NNN_*.up.sql` / `NNN_*.down.sql` pair into [migrations/](migrations/).
- **Tests**: none exist in this repo. Don't waste time looking for `_test.go` files.
- **Lint**: only `vue-tsc` via the frontend build. No Go linter is configured.

## Architecture

### Three binaries, one shared `internal/`

- [cmd/bot](cmd/bot/) — Telegram long-polling bot (not webhook-based).
- [cmd/admin](cmd/admin/) — admin HTTP API on `ADMIN_PORT` (default 8080), gated by `ADMIN_API_KEY`.
- [cmd/webapp](cmd/webapp/) — serves the Vue SPA from `./static` plus `/app/api/*` JSON endpoints (default port 8081). This is the binary the Mini App talks to.

All three load config, open a `pgx` pool, run migrations, then wire repositories → services → handlers.

### Layering (strict)

- [internal/repository/](internal/repository/) — `interfaces.go` is the canonical contract for every repo; concrete implementations live in per-entity `*_repo.go` files using raw `pgx`.
- [internal/service/](internal/service/) — business logic, constructed with repos and other services. `dashboardService` and `recommendationService` compose lower-level services rather than touching repos directly.
- [internal/handler/](internal/handler/) — thin transport layer split into `bot/`, `admin/`, `webapp/` sub-packages, one per binary.

### Auth flow (webapp)

Client sends Telegram `initData` in the `X-Telegram-Init-Data` header on the first request. The server validates the HMAC against the bot token, issues a Bearer token, and the client caches it in `localStorage`. Subsequent requests use `Authorization: Bearer …`. On a 401 the client retries with `initData` — see [web/src/api.ts](web/src/api.ts).

### Roles

`users.role` is `client` or `admin` (migration [021_user_role](migrations/021_user_role.up.sql)). Admin endpoints live under `/app/api/admin/*` and check the role server-side. There is no UI for promoting users — admins are set manually in the DB.

### Payments

[internal/payment/](internal/payment/) exposes a `Provider` interface. [cmd/webapp/main.go](cmd/webapp/main.go) picks `RobokassaProvider` when `ROBOKASSA_MERCHANT_LOGIN` is set, otherwise `DummyProvider` (instant success, for local dev). Callback verification flows through the `CallbackVerifier` interface.

### Frontend

Vue 3 + Vite + TypeScript. No Pinia or Vuex — state is local component state plus direct API calls. Pages live in [web/src/views/](web/src/views/); admin pages in [web/src/views/admin/](web/src/views/admin/). All endpoints flow through [web/src/api.ts](web/src/api.ts).

## Required env vars

See [.env.example](.env.example). The non-optional ones:

- `TELEGRAM_BOT_TOKEN` — required by both `bot` **and** `webapp` (webapp uses it to validate `initData`).
- `ADMIN_API_KEY` — required by the `admin` binary.
- `WEBAPP_URL` — used by the bot to launch the Mini App.
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`.
- Robokassa vars (`ROBOKASSA_MERCHANT_LOGIN`, `ROBOKASSA_PASSWORD1`, `ROBOKASSA_PASSWORD2`, `ROBOKASSA_IS_TEST`) are optional; leaving them empty selects the Dummy provider.

## Deployment

Production runs on **`94.247.128.101`** (referred to as the **"101 host"** in conversation). The Mini App is reachable at `https://fitness-bot.rnm.dev` (see `WEBAPP_URL` in [.env](.env)). Deploy mechanics (SSH user, deploy path, restart/reload command) are not currently scripted in the repo — confirm with the operator before pushing changes.

> **Shared host — handle with care.** The 101 host runs other developers' projects alongside fitness-bot. **Never** run `docker system prune`, `docker image/volume/network prune`, unscoped `docker stop/rm` patterns, or a `docker-compose down` from an unverified directory. Restart only the fitness-bot services (`docker-compose restart <svc>` or `docker-compose up -d --no-deps <svc>`). When in doubt, ask before running.

## Gotchas

- **Go version mismatch is intentional**: [go.mod](go.mod) declares `go 1.18` (minimum), the [Dockerfile](Dockerfile) builds with `golang:1.23-alpine`. Don't "align" them without confirming.
- **Postgres host port is 5444, not 5432.** [docker-compose.yml](docker-compose.yml) maps `5444:5432`. [.env.example](.env.example) lists `DB_PORT=5432` because that's the *container* port — for running a service from your host against the dockerized DB, use `5444`.
- **The webapp binary serves [web/dist/](web/dist/) from `./static`.** The Dockerfile's webapp stage copies `web/dist` over from the frontend build stage. If you `go run ./cmd/webapp` locally without first running `npm run build`, the SPA won't serve — only the API will.
- **Pre-built `bot`, `admin`, `webapp` binaries in the repo root** are leftovers from older builds. Ignore them; [cmd/](cmd/) is the source of truth.
