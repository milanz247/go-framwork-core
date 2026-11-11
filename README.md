# IT Support workspace

A minimal Go MVC application with verified registration, sign-in, sign-out and a
private dashboard. Built with Go templates, HTMX and SQLite.

## Development

Install Go 1.26 and Node.js 22.16 or newer, then run:

```sh
make setup
make dev
```

Set `LEGACY_VERIFY_URL` in `.env` to your full Laravel
`/api/verify-legacy-user` endpoint. Open http://localhost:8001.
The development launcher reads `.env` automatically.

## Registration

Existing IT Support email/password credentials are verified server-side by the
Laravel endpoint. A local account is created only when the response has
`success: true`, `password_verified: true`, a positive `legacy_user_id`, a name
and the same email. The name and email come from the verified response.
Passwords are stored locally as bcrypt hashes. Subsequent logins use the local
account. Missing configuration, failed verification and upstream failures never
create an account.

## Production

```sh
make build
set -a
. ./.env
set +a
make run
```

Production startup reads exported environment variables. Alternatively,
`docker compose up --build -d` reads `.env` and persists SQLite in `mvc-data`.
Set `PUBLIC_ORIGIN` to the browser origin and `COOKIE_SECURE=1` when using HTTPS.
Registration requires the configured Laravel endpoint, using HTTPS in production.

| Setting | Purpose |
| --- | --- |
| `PORT` | HTTP port, default 8001 |
| `APP_PORT` | Docker host port, default 8001 |
| `DATABASE_PATH` | SQLite file, default data/ticket-sync.db |
| `APP_ENV` | development enables live reload endpoints |
| `PUBLIC_ORIGIN` | Allowed browser origin |
| `COOKIE_SECURE` | 1 for HTTPS session cookies |
| `LEGACY_VERIFY_URL` | Full Laravel credential verification endpoint |

## Routes

HTML: `/login`, `/register`, `/dashboard`, POST `/logout`.
JSON: POST `/api/login`, POST `/api/register`, GET `/api/state`, POST `/api/logout`.
Authenticated JSON mutations require the session CSRF token in `X-CSRF-Token`.
`/api/state` returns only the current user's id, name, email and CSRF token.
`/healthz` is the health check. Other feature routes return 404.

The app creates only users and sessions tables. Existing databases and user
accounts remain compatible; previous feature tables are left untouched and are
no longer read or updated. There are no reporting integrations, background jobs,
SMTP, Google credentials or Jira connections.

## Checks

```sh
make check
make build
```

## Contributing

See [the development guide](docs/development.md) for the folder structure,
request flow and where to make changes. Run `make help` for all commands.
Use `make fmt` followed by `make check` before sharing changes. Editable frontend
files live in `resources/`; compiled files in `public/assets/` are embedded into Go.

## Source structure

- `app/`: application authentication, handlers, repositories and view composition.
- `framework/`: reusable HTTP, database, rendering, security and runtime helpers and development tooling.
- `app/app.go`: wires configuration, database and routes to the framework core.
- `config/`: application configuration validation.
- `routes/`: HTTP route registration.
- `resources/`: editable views, CSS, JavaScript and icons.
- `public/assets/`: built browser assets embedded in the executable.
- `cmd/server/`: executable entry point.

Framework code never imports application packages. See the development guide for
component conventions, request flow and build commands.
