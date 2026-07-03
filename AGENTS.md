# AGENTS.md — Kasseapparat

Compact cheat-sheet for OpenCode. Omit anything an agent could guess from filenames.

---

## Repo layout

- **Go backend**: `backend/` (module `github.com/potibm/kasseapparat`)
- **React frontend**: `frontend/` (Vite + React Admin + Tailwind v4)
- **Infra**: `docker-compose.yaml` (OpenObserve + OTel collector + Redis + Mailhog)
- **No `go.work` at root**: Run Go commands from `backend/` or use `mise` tasks.

Frontend apps (two SPAs routed by React Router):

- `/` → POS (Point of Sale) — protected, requires auth
- `/admin/*` → React Admin dashboard — public

Path aliases (vite + tsconfig): `@core`, `@admin`, `@pos`.

---

## Everyday commands (mostly via `mise`)

Install mise tools once: `mise install`  
Full setup (deps + infra): `mise run setup`  
Dev (hot-reload both): `overmind s --timeout 10` (uses Procfile)

Backend only: `mise run be:dev` (Air, port 3001)  
Frontend only: `mise run fe:dev` (Vite, port 3000, HTTPS, proxies `/api` → :3001)

Test everything: `mise run test`  
Backend tests: `mise run be:test`  
Frontend tests: `mise run fe:test`

Lint everything: `mise run lint`  
Lint with auto-fix: `mise run lint --fix`

Docker image check: `mise run docker:build` (uses a custom `kasseapparat-builder` buildx builder)

---

## Critical gotchas

### `cmd/assets` must exist before any Go build

`backend/cmd/serve.go` embeds `//go:embed assets`. If the directory is missing, **compilation fails**.

- `mise run be:setup` creates it (plus a dummy `index.html`) and copies `.env.example` → `.env`.
- `mise run be:lint` also creates a dummy file for this reason.
- `mise run be:test` depends on `be:setup`, so tests via `mise` are safe; running `go test ./...` directly without `be:setup` first will fail.
- Dockerfile copies the real frontend build into `backend/cmd/assets`.

### Config loading order

1. `backend/config/config.yaml` (committed defaults)
2. `backend/config/config.local.yaml` (gitignored overrides, merged if present)
3. `.env` (loaded by `godotenv`)
4. Environment variables (`APP_LOG_LEVEL` maps to `app.log_level`)
5. CLI flags (`--log-level`, `--port`, etc.)

Use `config/config.local.yaml` for local secrets; do not edit `config.yaml`.
Generate a fresh config with: `go run . config create`

### Frontend dev proxy

Vite proxies `/api` to `http://127.0.0.1:3001`.  
The backend dev server (`air`) runs on **3001**, not 3000.

### Frontend tests need `--no-webstorage`

The `fe:test` task sets `NODE_OPTIONS="--no-webstorage"`. Running `npm run test` directly may behave differently.

### Auth is JWT-based (not OIDC)

The backend uses `appleboy/gin-jwt/v3` for username/password auth with JWT tokens. The admin panel has its own auth flow via `auth-provider.ts`. Do not assume OIDC/JWKS patterns.

---

## Lint / typecheck / test pipeline

CI order: `deps:install` → `lint` → `test` → `docker:build`.

- **Backend lint**: `golangci-lint` (config in `backend/.golangci.yaml`). Uses `gofumpt` + `golines` (120 cols).
- **Frontend lint**: ESLint + Prettier + `tsc --noEmit` + `dotenv-linter`.
- **Repo lint**: Prettier over `*`, `.github`, and `backend/**/*.{json,yml,yaml,md}`.

---

## Testing

- **Backend**: `go test -v -coverprofile=coverage.out ./...` (run from `backend/`).
- **Frontend**: `vitest`, jsdom, globals enabled, coverage via v8 (`frontend/coverage/lcov.info`).
- SonarCloud ingests `backend/coverage.out` and `frontend/coverage/lcov.info`.

---

## Database / data directories

SQLite file lives in `backend/data/`.  
The `be:setup` task creates the `data/` directory if missing.

---

## Releasing

- Semantic Release with Angular commit preset.
- Docker release triggers on tags matching `[0-9]+.[0-9]+.[0-9]+`.
- Image is signed with Cosign and attested with SBOM.
- PR titles are validated by `amannn/action-semantic-pull-request`.

---

## Infra (local)

`docker compose up -d` starts:

- Mailhog: http://localhost:8025 (SMTP: localhost:1025)
- OpenObserve UI: http://localhost:8027 (admin@example.com / password123)
- OTel gRPC: localhost:4317
- Redis: localhost:6379
- RedisInsight: http://localhost:8026

Backend dev (`air`) defaults to `--otel-endpoint=localhost:4317`.

---

## Style notes

- Go: `slog` only (depguard blocks `logrus`). Use snake_case for slog keys.
- Go: avoid `math/rand` (use `math/rand/v2`).
- TS: `no-console` and `no-alert` are errors. `_` prefix ignores unused vars.
- TS: new `.js` files are forbidden by ESLint (use `.ts`/`.tsx`).
