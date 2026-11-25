![Go CI](https://img.shields.io/github/actions/workflow/status/ozaitsev92/go-react-todo-list/go.yml?branch=main&label=go)
![React CI](https://img.shields.io/github/actions/workflow/status/ozaitsev92/go-react-todo-list/react.yml?branch=main&label=react)

# Go React TODO App

Simple TODO app built with Go, React and MongoDB.

<img src="./img/todo.png">

### Usage

Run `make start`, then open http://localhost:8081. On first run it generates a gitignored `.env` with random secrets (delete it to rotate). Compose serves plain HTTP, so it sets `JWT_SECURE_COOKIE=false`.

---

## Frontend

React 19, React Router 7, React Bootstrap, built with [Vite](https://vite.dev/). The source is in [`frontend/`](frontend). Node 22.

```bash
cd frontend
npm ci
npm run dev        # http://localhost:8081, proxies /api to :8080
npm test           # no backend needed
npm run lint
npm run format
npm run build      # into dist/
```

The app calls the API at `/api` on its own origin, and the Vite dev server, nginx or the k8s ingress routes it to the backend. So one image works everywhere and no CORS is needed. The session lives only in the backend's HttpOnly cookie.

The production image is unprivileged nginx with caching, gzip, a strict CSP, `/healthz` and the `/api` proxy (`API_UPSTREAM`).

---

## Backend

Go 1.26, [Gin](https://gin-gonic.com/), MongoDB. The source is in [`backend/`](backend).

### Running

```bash
cd backend
make run                # in-memory storage, dev-only key, http://localhost:8080
make test               # unit tests, no Docker needed
make test-integration   # + MongoDB integration tests (Docker, or set MONGO_TEST_URI)
make lint vuln          # go vet + golangci-lint, govulncheck
```

### Configuration

Defaults < optional TOML file (`CONFIG_PATH`, see [`config.example.toml`](backend/config/config.example.toml)) < env vars. Invalid config stops startup.

| Env var | Default | Notes |
|---|---|---|
| `JWT_SIGNING_KEY` | (none, **required**) | ≥ 32 bytes, e.g. `openssl rand -hex 32`. Never commit it. |
| `STORAGE` | `mongo` | `memory` keeps data in-process only (lost on restart) |
| `MONGO_URL` / `MONGO_DB_NAME` | `mongodb://localhost:27017` / `todo` | |
| `BIND_ADDR` | `8080` | port |
| `JWT_SESSION_LENGTH` | `30` | minutes |
| `JWT_SECURE_COOKIE` | `true` | set `false` only for plain-HTTP development |
| `JWT_COOKIE_DOMAIN`, `JWT_ISSUER` | `""`, `todo-api` | |
| `PASSWORD_HASH_COST` | `12` | bcrypt cost, 10–16 |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `MAX_BODY_BYTES` | `65536` | larger requests get `413` |
| `MONGO_TIMEOUT`, `READ_HEADER_TIMEOUT`, `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT`, `GRACEFUL_TIMEOUT` | `10`, `5`, `15`, `15`, `60`, `15` | seconds |

### API

Authentication uses an `HttpOnly`, `SameSite=Lax` cookie named `jwt-token`, set by `POST /v1/users/login`. Errors have the shape `{"error": "message"}`.

| Method & path | Success | Errors |
|---|---|---|
| `POST /v1/users` `{"email","password"}` | `201` user | `400` bad JSON, `409` email taken, `422` invalid email / password not 8–72 bytes |
| `POST /v1/users/login` `{"email","password"}` | `200` user + cookie | `401` invalid email or password |
| `POST /v1/users/logout` | `204`, cookie cleared | |
| `GET /v1/users/current` | `200` user | `401` |
| `GET /v1/tasks` | `200` own tasks, oldest first | `401` |
| `POST /v1/tasks` `{"text"}` | `201` task + `Location` | `400`, `413`, `422` text not 1–1000 chars |
| `GET /v1/tasks/{id}` | `200` task | `400` invalid id, `404` |
| `PUT /v1/tasks/{id}` `{"text"}` | `200` task | `400`, `404`, `422` |
| `PUT /v1/tasks/{id}/mark-completed`, `/mark-not-completed` | `200` task | `400`, `404` |
| `DELETE /v1/tasks/{id}` | `204` | `400`, `404` |
| `GET /healthz` / `GET /readyz` / `GET /metrics` | liveness / readiness (pings MongoDB) / Prometheus | `503` from `/readyz` when the DB is unreachable |

Another user's task returns the same `404` as a missing one.

### Architecture

```
cmd/app ─▶ internal/app (composition root: config, storage, wiring, lifecycle)
              │
              ├─▶ internal/controller/http/v1   Gin handlers, DTOs, JWT cookie, middleware
              │        │
              ├─▶ internal/usecase              ownership checks, registration, authentication
              │        │                        (PasswordHasher port)
              │        ▼
              │   internal/domain/{task,user}   entities, invariants, repository ports
              │        ▲
              └─▶ internal/infrastructure       mongo + memory repositories, bcrypt hasher
pkg/{httpserver,logger,mongodb}                 generic helpers with no app imports
```

Both repository implementations pass one shared contract suite, so `STORAGE=memory|mongo` are interchangeable.

### Design decisions and trade-offs

- **Stateless JWT in an HttpOnly cookie.** XSS can't read the token. CSRF is mitigated by `SameSite=Lax` and the same-origin deployment (the API sends no CORS headers). Tokens can't be revoked before they expire (30 min by default); refresh tokens or a `jti` denylist would fix that.
- **Uniqueness comes from the database.** Emails are normalized (trimmed, lowercased), and a unique MongoDB index enforces one account per email, even under concurrent sign-ups.
- **Not implemented yet:** rate limiting on the auth endpoints, pagination on `GET /v1/tasks`, and an OpenAPI spec.

### Kubernetes

The manifests in [`k8s/`](k8s) expect two Secrets that are not committed:

```bash
kubectl -n todo-app create secret generic mongo-secrets \
  --from-literal=username=todo --from-literal=password="$(openssl rand -hex 24)"
kubectl -n todo-app create secret generic backend-secrets \
  --from-literal=jwt-signing-key="$(openssl rand -hex 32)" \
  --from-literal=mongo-url="mongodb://todo:<password>@mongo-0.mongo:27017/?authSource=admin"
```

The ingress serves the API under `/api` and must use HTTPS, because the session cookie is `Secure` (set `JWT_SECURE_COOKIE=false` for a plain-HTTP local cluster).
