# JumpAlias

## What it is

JumpAlias is a small internal tool for mapping short aliases to full URLs.
Teams register an alias once, then jump to the destination through a redirect path.

## What you can do

- Register a new alias → destination mapping
- Browse every registered shortcut
- Follow `/j/:alias` to land on the destination
- Reject invalid aliases and non-http(s) destinations
- Block duplicate aliases
- Persist data in SQLite

## Stack

| Layer | Tools |
|-------|--------|
| API | Go, Gin, GORM, SQLite |
| UI | React, TypeScript, Vite, Axios |

## Layout

```text
cmd/server/           process entry (bootstraps the app)
internal/
  di/                 settings + dependency wiring
  router/             HTTP engine and route mount
  handler/            request/response adapters
  service/            alias business rules
  repository/         SQLite persistence
  model/              Shortcut entity
  validation/         alias + destination checks
frontend/             React UI
```

Flow:

```text
ShortcutAPI → AliasManager → ShortcutStore → SQLite
         ↑
      di.Bootstrap
         ↓
     BuildEngine
```

## Requirements

- Go 1.23+
- Node.js 20+
- npm

## Start the API

```bash
go mod tidy
go run ./cmd/server
```

Listens on [http://localhost:8080](http://localhost:8080).

Environment knobs:

- `LISTEN_PORT` — listen port (default `8080`)
- `SQLITE_FILE` — SQLite file path (default `jumpalias.db`)

## Start the UI

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

UI: [http://localhost:5173](http://localhost:5173).

Point the UI at the API with `VITE_BACKEND_ORIGIN` (default `http://localhost:8080`).

## HTTP surface

| Method | Path | Behavior |
|--------|------|----------|
| GET | `/ready` | Liveness probe |
| POST | `/api/shortcuts` | Register a shortcut (`alias`, `destination`) |
| GET | `/api/shortcuts` | List shortcuts |
| GET | `/j/:alias` | Redirect to destination |

Example create body:

```json
{
  "alias": "design-system",
  "destination": "https://example.com/design"
}
```

Successful create responses include `id`, `alias`, `destination`, and `registeredAt`.

## Notes

- Aliases are unique and stored lowercase.
- Destinations must be absolute `http` or `https` URLs.
- Auth is out of scope for this version.
- Local UI expects a local API process.

## Later ideas

- Search / filter aliases
- Edit or remove shortcuts
- Pagination
- Authn/authz
- Container packaging
- Broader automated coverage
