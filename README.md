# labbaika

Monorepo with a Go backend (Echo + Cobra, GORM/Postgres) and a React (CRA + MUI) frontend.

## Structure

```
backend/    Go services: payslip-svc, cashbook-svc (see backend/readme.md)
frontend/   React (TypeScript, MUI) client
```

## Services (backend)

Single Go binary, subcommand per service (Cobra):

| Command | Description |
|---|---|
| `payslip-svc` | REST API for payslip generation (docx template based) |
| `cashbook-svc` | REST API for cashbook (money in/out) |

Each runs an Echo server on `:8080` under `/api/v1`, with trace ID, JSON access logging, and error-handling middleware.

## Dependencies

| | Version | Notes |
|---|---|---|
| Go | 1.20 | `backend/go.mod` |
| Node | 16.18.1 | pinned in `frontend/.nvmrc`, use `nvm use` |
| Postgres | 14+ (any recent) | only DB in use; `docker-compose.yaml` ships `postgres:latest` |
| Docker + docker-compose | — | only needed for the `docker-compose up` path |

`cashbook-svc` additionally needs the DB connection enabled — `initDatabase()` is commented out in `backend/cmd/app/root.go`'s `cobra.OnInitialize` — and its table created: see `backend/pkg/database/migration/`.

## Run locally

Everything (Postgres + backend + frontend):
```
docker-compose up
```

Backend only:
```
cd backend
# fill backend/.env: PG_HOST/PG_PORT/PG_USER/PG_PASSWORD/PG_DATABASE (+ PG_RO_* for read replica, can mirror the above)
go run cmd/main.go payslip-svc   # or cashbook-svc
```
Env vars are loaded from `backend/.env` via `godotenv` when `APP_ENV` is empty/`dev`/`local` (see `backend/cmd/config/env.go`).

Frontend only:
```
cd frontend
nvm use          # Node 16.18.1
npm install
npm start
```

## Build

```
cd backend && make build   # or: go build -o labbaika ./cmd
```

## Details

- Backend layout, folder conventions, migrations: [backend/readme.md](backend/readme.md)
- Frontend: [frontend/README.md](frontend/README.md) (CRA defaults)
