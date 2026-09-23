# Go Links

## Overview

Go Links is a simple internal URL shortener that allows users to create memorable shortcuts for frequently used URLs. Users can create shortcuts, browse existing links, and access the destination URL through a redirect endpoint.

## Features

- Create a Go Link
- List all Go Links
- Redirect using shortcut
- Input validation
- Duplicate slug detection
- SQLite persistence

## Tech Stack

**Backend**
- Go
- Gin
- GORM
- SQLite

**Frontend**
- React
- TypeScript
- Vite
- Axios

## Project Structure

```text
cmd/server/              Application entrypoint (starts DI)
internal/
  di/                    Single DI composition root (wires dependencies)
  router/                HTTP routing and middleware configuration
  handler/               HTTP adapters (Gin handlers + response DTOs)
  service/               Business rules (use-cases)
  repository/            Persistence port + GORM adapter
  model/                 Domain entities
  validation/            Input normalization and validation
frontend/                React + TypeScript UI
README.md
```

Dependency flow:

```text
Handler → Service → Repository → Database
                ↑
            di wires these
                ↓
             router
```

## Prerequisites

- Go 1.23+
- Node.js 20+
- npm

## Running the Backend

```bash
go mod tidy
go run ./cmd/server
```

Runs on [http://localhost:8080](http://localhost:8080).

Optional environment variables:

- `PORT` — server port (default `8080`)
- `DB_PATH` — SQLite file path (default `golinks.db`)

## Running the Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Runs on [http://localhost:5173](http://localhost:5173).

The frontend API base URL is configured with `VITE_API_BASE_URL` (defaults to `http://localhost:8080`).

## API Endpoints

| Method | Endpoint     | Description                 |
|--------|--------------|-----------------------------|
| POST   | `/api/links` | Create a Go Link            |
| GET    | `/api/links` | List Go Links               |
| GET    | `/go/:slug`  | Redirect to destination URL |

## Assumptions

- Slugs are unique.
- URLs must be valid (`http` or `https`).
- SQLite is used for simplicity.
- Authentication is intentionally omitted.
- The frontend communicates with the backend running locally.

## Future Improvements

- Search and filter shortcuts
- Edit and delete existing links
- Pagination for large datasets
- Authentication and authorization
- Docker support
- Unit and integration tests
- Configurable backend URL via environment variables (frontend already uses `VITE_API_BASE_URL`)
