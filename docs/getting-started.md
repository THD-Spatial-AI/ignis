# Getting started

ignis is an HTTP API that computes TABULA heating demand for building archetypes. It ships as three containers: a Caddy reverse proxy, the Go application, and its own PostgreSQL instance.

Three setup paths are available. Pick one:

| Path | Use when | Needs |
|---|---|---|
| [Quick start](#quick-start) | Trying the API out | Docker only |
| [Docker from source](#docker-from-source) | Developing ignis itself | Docker, `caddy` CLI, source checkout |
| [Manual setup](#manual-setup) | Running without containers | Go, PostgreSQL, source checkout |

## Configuration files

All three paths read the same two environment files. They are separate on purpose, because they are consumed at different points.

| File | Read by | Holds |
|---|---|---|
| `.env` | Docker Compose, on the host | `APP_PORT`, `HOST_HTTPS_PORT`, `CADDY_DATA_DIR`, `IGNIS_IMAGE_TAG` |
| `docker.env` | Passed into the containers | `DB_*`, `POSTGRES_*`, `ALLOWED_ORIGINS`, `IGNIS_API_KEY` |

Both live in `environment/`. `docker.env` is committed with local development defaults, including a placeholder API key and database password.

!!! danger "Replace the committed defaults before any deployment"
    `docker.env` contains `IGNIS_API_KEY=supersecret123`, `POSTGRES_PASSWORD=postgres`, and `DB_SSL_MODE=disable`. These are development placeholders. See [Deployment](#deployment).

## Authentication

ignis has no authentication of its own. The reverse proxy gates every request behind an `X-Api-Key` header, matched against `IGNIS_API_KEY` from `docker.env`:

```bash
curl -k -H "X-Api-Key: supersecret123" https://localhost/some-endpoint
```

Requests without a valid key receive `403 Forbidden` from the proxy and never reach the app.

`/health` is the single exemption, so container health checks and uptime monitors need no credential. It reports process liveness only and does not check the database connection.

!!! danger "Never expose ignis directly"
    Because authentication lives entirely in the proxy, ignis itself must never be reachable except through it. None of the compose files publish a port for the app, and that should not change.

---

## Quick start

For trying out the API with nothing on the host but Docker. Pulls pre-built images from GHCR, so no Go toolchain, no `caddy` install, and no `.env` file are needed.

### 1. Get the files

You need three files in one directory: `docker-compose.quickstart.yml`, `Caddyfile`, and `docker.env`. Cloning the repository is the simplest way to get them:

```bash
git clone https://github.com/thd-spatial-ai/ignis.git
cd ignis/environment
```

### 2. Start the stack

```bash
docker compose -f docker-compose.quickstart.yml up -d
```

This starts `ignis-db`, then `ignis-app` once the database reports healthy, then `ignis-reverse-proxy` once the app reports healthy. Only the proxy publishes a host port (`HOST_HTTPS_PORT`, default `443`).

### 3. Seed the database

!!! warning "Required before first use, and destructive"
    A fresh `ignis-db` volume is empty. Seeding drops and recreates all country tables, so it is gated behind the `seed` profile and never runs automatically. Until it has run once, every endpoint that reads the schema will fail.

```bash
docker compose -f docker-compose.quickstart.yml --profile seed run --rm ignis-build-db
```

The TABULA workbook is baked into the `ignis-build-db` image, so there is nothing to download.

### 4. Verify

```bash
docker compose -f docker-compose.quickstart.yml exec ignis-db psql -U postgres -d ignis -c "\dt tabula.*"
curl -k -s -o /dev/null -w '%{http_code}\n' https://localhost/health
```

The first lists the seeded tables. The second returns `200`.

### Pinning a version

Both images default to the `latest` published release. To pin a specific one, set `IGNIS_IMAGE_TAG` in the environment before starting:

```bash
export IGNIS_IMAGE_TAG=0.2.4-alpha
docker compose -f docker-compose.quickstart.yml --profile seed pull
docker compose -f docker-compose.quickstart.yml up -d
```

!!! info "Tag format"
    Published image tags carry no `v` prefix, even though the Git tags do. Release `v0.2.4-alpha` publishes as `0.2.4-alpha`. Export the variable rather than prefixing a single command, so the app and seed images come from the same release.

### Certificate warning is expected

The quickstart file keeps Caddy's local CA in a Docker-managed volume rather than a host bind mount, so it is never added to your OS or browser trust store. Requests to `https://localhost` will show an untrusted-certificate warning.

Click through it in the browser, or pass `-k` (curl), `--no-check-certificate` (wget), or "disable SSL verification" (Postman).

!!! info "Using Swagger UI"
    Calling the API from Swagger UI's "Try it out" on this site's [API reference](api.md) needs one extra step. Browser JavaScript cannot click through a certificate warning the way a manual page load can, so a `fetch()` to an untrusted origin fails outright. Open `https://localhost` directly in a new tab first and click through the warning. Most browsers then trust that origin for the rest of the session.

### Tearing down

```bash
docker compose -f docker-compose.quickstart.yml down -v
```

The `-v` removes the database volume, so the next start needs seeding again. Omit it to keep the seeded data.

---

## Docker from source

Builds `ignis-app` and `ignis-build-db` from this checkout, so it is the only path that picks up local code changes.

### Prerequisites

Beyond Docker, this path needs the `caddy` CLI on the host, with `caddy trust` run once. Doing so installs a local CA into your OS and browser trust store, which the proxy then reuses, so `https://localhost` is trusted without a warning.

```bash
caddy trust
```

### 1. Configure

```bash
cd environment
cp .env.example .env
```

`.env` must define:

| Variable | Description | Default |
|---|---|---|
| `CADDY_DATA_DIR` | Host directory holding the `caddy trust` CA | required, no default |
| `APP_PORT` | The app's internal listen port | `8080` |
| `HOST_HTTPS_PORT` | The proxy's published port | `443` |

!!! warning "CADDY_DATA_DIR has no default on purpose"
    A wrong value silently produces an untrusted certificate, so a missing one fails validation instead. Compose reports it by name: `required variable CADDY_DATA_DIR is missing a value`.

### 2. Start and seed

```bash
docker compose up -d
docker compose --profile seed run --rm ignis-build-db
```

The seed step is destructive in the same way as the quickstart path.

### 3. Verify

```bash
docker compose exec ignis-db psql -U postgres -d ignis -c "\dt tabula.*"
curl -s -o /dev/null -w '%{http_code}\n' https://localhost/health
```

No `-k` is needed here, since the certificate chains to the CA you trusted.

### Tearing down

```bash
docker compose down -v
```

---

## Manual setup

For local development without containers.

### Prerequisites

| Dependency | Requirement |
|---|---|
| Go | 1.26 or later |
| PostgreSQL | 15 to 17 |
| Git LFS | required for the TABULA workbook |

!!! warning "Install Git LFS before cloning"
    `data/tabula-calculator-lite.xlsx` is stored via Git LFS. Without `git lfs install`, the checkout contains a small text pointer instead of the workbook, and `build_db` fails with `zip: not a valid zip file`.

    ```bash
    git lfs install
    git clone https://github.com/thd-spatial-ai/ignis.git
    cd ignis
    ```

    If you have already cloned without it, run `git lfs install && git lfs pull`. Verify with `ls -l data/tabula-calculator-lite.xlsx`, which should be roughly 11 MB rather than a few hundred bytes.

The workbook is committed to the repository, so there is nothing to download from [episcope.eu](https://episcope.eu/iee-project/tabula/).

### Configuration

The manual path reads `environment/docker.env` for its database and CORS settings. Copy it if you want local values, and point `DB_HOST` at your own PostgreSQL instance:

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `ignis-db` (set to `localhost` for manual setup) |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | required, no default |
| `DB_NAME` | Database name | `ignis` |
| `DB_SSL_MODE` | TLS mode (`require` or `disable`) | `disable` in the committed file |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins | see below |

!!! info "ALLOWED_ORIGINS"
    Only needed when a browser-based client calls ignis directly. For server-to-server calls, the intended deployment model, it can be left unset.

!!! warning "DB_SSL_MODE"
    Use `disable` only for local development against a database on the same machine. Set `require` in all other environments.

!!! warning "No API key on this path"
    Running the binary directly bypasses the reverse proxy, so there is no `X-Api-Key` gate and no CORS preflight handling. Bind it to localhost only.

### Build

```bash
go build -buildvcs=false -o bin/ ./cmd/...
```

This produces `bin/ignis`, `bin/build_db`, and `bin/validate`.

### Load the database

```bash
./bin/build_db
```

Destructive: drops and recreates the `tabula` schema.

### Run the API

```bash
./bin/ignis
```

Starts on `APP_PORT`, default `8080`.

### Validate

Runs the full 17-level TABULA calculation pipeline against every row in the database and checks that each result stays within 2% of the reference value from the workbook.

```bash
./bin/validate
```

This path only: there is no containerised `validate`. See the [validation report](validation.md) for current results.

---

## Deployment

Use `docker-compose.prod.yml`, which pulls published images and needs no source tree on the target machine. Copy four files across: the compose file, `.env`, `docker.env`, and `Caddyfile`.

### 1. Prepare `docker.env`

!!! danger "Do not deploy the committed docker.env"
    Change all of the following before starting the stack:

    - `POSTGRES_PASSWORD` and `DB_PASSWORD` to a real credential
    - `DB_SSL_MODE` to `require`
    - `IGNIS_API_KEY` to a rotated value
    - `ALLOWED_ORIGINS` to the real browser origins, or unset for server-to-server only

### 2. Prepare `.env`

`CADDY_DATA_DIR` is required, `APP_PORT` and `HOST_HTTPS_PORT` default to `8080` and `443`. Set `IGNIS_IMAGE_TAG` to pin a release, which is strongly advised for anything you do not want moving underneath you.

### 3. Set the site address

!!! warning "The Caddyfile ships with `localhost` as its site address"
    Change it to the deployment's real domain, or Caddy will neither serve it nor provision a certificate for it.

### 4. Pull, start, seed

```bash
docker compose -f docker-compose.prod.yml --profile seed pull
docker compose -f docker-compose.prod.yml up -d
docker compose -f docker-compose.prod.yml --profile seed run --rm ignis-build-db
```

!!! warning "Pull first, and include the seed profile"
    With `IGNIS_IMAGE_TAG` unset, a host that already has `latest` cached keeps running the old build after a release, with nothing on that host revealing it. The `--profile seed` on the pull is what fetches `ignis-build-db`, since profile-gated services are otherwise skipped.

Seed only on first deployment. Running it against a populated database drops every country table.

### 5. Verify

```bash
curl -s -o /dev/null -w '%{http_code}\n' https://your-domain/health
curl -s -o /dev/null -w '%{http_code}\n' https://your-domain/
curl -s -o /dev/null -w '%{http_code}\n' -H "X-Api-Key: your-key" https://your-domain/
```

Expect `200`, `403`, then a response from the app. A `200` on the second call means the API key gate is not working, and the deployment should be stopped.

---

## Health checks

Each container reports health to Compose, which is what orders startup: `ignis-app` waits for a healthy `ignis-db`, and `ignis-reverse-proxy` waits for a healthy `ignis-app`.

The app is probed every 30 seconds, with a 10 second start period during which failures do not count against the retry budget. Successful probes are not logged, so `docker compose logs ignis-app` shows real traffic rather than health check noise. Failing probes are logged.

`/health` reports process liveness and does not test the database connection, so a healthy container does not by itself mean the schema is seeded or reachable.
