# Development guide

## Structure

```text
app/
  app.go                Application startup wiring
  auth/                 Registration, legacy verification and sessions
  handlers/             HTML pages and JSON API
  repositories/         User/session queries and application schema
  views/                Application layout selection and view helpers
cmd/server/             Executable entry point
config/                 Application configuration validation
framework/
  core/                 HTTP lifecycle and graceful shutdown
  database/             SQLite connection and query helpers
  dev/                  Development reload endpoints
  env/                  Environment access
  http/                 Middleware and JSON response helpers
  security/             Password verification
  support/              Shared values, tokens and timestamps
  view/                 Template loading, helpers and rendering
  tooling/              Asset build, dev launcher, formatter and Air config
resources/
  views/
    layouts/            Base, app and authentication layouts
    components/app/     Sidebar, navbar, footer and account controls
    components/auth/    Shared form, logo and footer
    components/ui/      Buttons, icons, fields and alerts
    pages/              Login, register and dashboard
  css/                  Editable styles
  js/                   Browser behavior
  icons.svg             Used icon symbols
public/assets/          Built assets embedded in the executable
routes/                 HTTP route registration
```

`data/` holds runtime SQLite files. `private/` contains preserved local private
files. Both are ignored by Git and Docker. `bin/`, `.tmp/` and `node_modules/` are
local build/tool outputs.

## Workflow

Use Go 1.26+ and Node.js 22.16+. Run commands from the project root:

```sh
make setup      # Install dependencies and create .env if missing
make dev        # Run with server and browser reload
make fmt        # Format Go and HTML
make check      # Check formatting, run go vet and build assets
make build      # Build bin/ticket-sync
```

Set `LEGACY_VERIFY_URL` in `.env`. Development and Docker Compose load `.env`;
production startup reads exported variables. Restart development after changing
`.env`. `make help` lists commands. `make clean` removes `bin/` and `.tmp/`, keeping
runtime data; reinstall Air with `make install` afterward.

Edit frontend sources in `resources/`. `make assets` compiles/copies them into
`public/assets/`, which is embedded and served at `/static/`. Direct changes to
built assets are overwritten.

## Layouts and components

Each HTML file owns one namespaced template definition:

- `layouts/base.html` owns the document, head and scripts.
- `layouts/app.html` composes sidebar, navbar, dashboard and footer.
- `layouts/auth.html` composes authentication branding and page content.
- Application components receive the page context (`.`).
- UI components receive explicit values through `dict`.
- Login/register pages configure the shared `auth/form` component.

```gotemplate
{{template "ui/button" (dict "Type" "submit" "Class" "primary" "Label" "Save")}}
{{template "ui/icon" "overview"}}
{{template "ui/avatar" .User.name}}
```

All nested HTML files in `resources/views/` are embedded automatically. Add a
unique `{{define "namespace/name"}}` and include it with `template`; no file list
needs updating. Keep `#workspace`, CSRF fields and `data-action` hooks intact.

Neutral light/dark theme tokens live in `resources/css/app.css`. The interface is
rendered by Go templates and HTMX with a Laravel starter-kit-inspired layout.
Mobile navigation and the account dialog are controlled in `resources/js/app.js`.

## Framework boundary

Framework packages do not import application code. Domain logic and schema
belong in `app/`; route choices belong in `routes/`; startup wiring belongs in
`app/app.go`.

Requests pass through middleware and routes to handlers. Handlers call auth and
repository functions, then render HTML or JSON. Repositories own domain SQL;
`framework/database` owns the connection and generic queries. `app/views` selects
layouts, while `framework/view` loads and renders the supplied template filesystem
and provides the shared `initial` and `dict` template helpers.
Existing database files remain compatible and old tables are not deleted.

## Core and tooling

`cmd/server/main.go` calls `app.Run()`. The application validates its settings,
opens its repositories and passes its HTTP routes to `framework/core.Run()`.
The core owns server timeouts, process signals and graceful shutdown; it does not
import application packages.

Development helpers live in `framework/dev/`, and routes supply the directories
they should watch. Build scripts and Air configuration live in
`framework/tooling/`. `make dev`, `make build` and `make check` remain the entry
points for working with them.
