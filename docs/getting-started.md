# Getting started

ignis is an HTTP API that computes TABULA heating demand for building archetypes. It ships as three containers: a Caddy reverse proxy, the Go application, and its own PostgreSQL instance.

Three setup paths are available. Pick one:

| Path | Use when | Needs |
|---|---|---|
| [Quick start](#quick-start) | Trying the API out | Docker only |
| [Docker from source](#docker-from-source) | Developing ignis itself | Docker, `caddy` CLI, source checkout |
| [Manual setup](#manual-setup) | Running without containers | Go, PostgreSQL, source checkout |

Both Docker paths share the same seed, verify and teardown steps; only the compose file differs.

| Path | Compose prefix |
|---|---|
| Quick start | `docker compose -f docker-compose.quickstart.yml` |
| Docker from source | `docker compose` |
| [Deployment](#deployment) | `docker compose -f docker-compose.prod.yml` |

## Configuration files

All paths read `.env`, interpolated on the host, plus four files under `environment/env/` passed into the containers. They are split so each service receives only the variables it reads: the internet-facing proxy never gets a database password.

| File | Read by | Holds |
|---|---|---|
| `.env` | Docker Compose, on the host | `APP_PORT`, `HOST_HTTPS_PORT`, `CADDY_DATA_DIR`, `IGNIS_IMAGE_TAG` |
| `env/common.env` | app, proxy | `ALLOWED_ORIGINS` |
| `env/db.env` | db | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `env/app.env` | app, build_db | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE` |
| `env/proxy.env` | proxy | `IGNIS_API_KEY`, `IGNIS_SITE_ADDRESS` |

The `env/` files are committed with local development defaults, including a placeholder API key and database password. Replace all of them before any deployment: see [Deployment](#deployment).

## Authentication

The reverse proxy gates every request behind an `X-Api-Key` header, matched against `IGNIS_API_KEY` from `env/proxy.env`. Requests without a valid key receive `403 Forbidden` and never reach the app.

```bash
curl -k -H "X-Api-Key: supersecret123" https://localhost/some-endpoint
```

`/health` is the single exemption, so container health checks and uptime monitors need no credential.

!!! danger "Never expose ignis directly"
    Authentication lives entirely in the proxy, so ignis itself must never be reachable except through it. None of the compose files publish a port for the app, and that should not change.

## Certificates

How `https://localhost` behaves depends on where Caddy's local CA is stored.

| Path | CA location | Result |
|---|---|---|
| Quick start | Docker-managed volume | Never enters your trust store. Expect an untrusted-certificate warning. |
| Docker from source | Host directory (`CADDY_DATA_DIR`), created by `caddy trust` | Trusted, no warning, no `-k` needed. |

On the quickstart path, click through the warning in the browser, or pass `-k` (curl), `--no-check-certificate` (wget), or "disable SSL verification" (Postman).

!!! info "Using Swagger UI"
    Calling the API from Swagger UI's "Try it out" on the [API reference](api.md) needs one extra step on the quickstart path. Browser JavaScript cannot click through a certificate warning the way a manual page load can, so a `fetch()` to an untrusted origin fails outright. Open `https://localhost` directly in a new tab first and click through the warning. Most browsers then trust that origin for the rest of the session.

---

## Quick start

For trying out the API with nothing on the host but Docker. Pulls pre-built images from GHCR, so no Go toolchain, no `caddy` install, and no `.env` file are needed.

### 1. Get the files

You need `docker-compose.quickstart.yml`, the `caddy/` directory, and the `env/` directory, all in one working directory. Cloning the repository is the simplest way to get them:

```bash
git clone https://github.com/thd-spatial-ai/ignis.git
cd ignis/environment
```

!!! warning "Docker Desktop: work from a directory under your home folder"
    Docker Desktop only shares paths under your home directory, or another folder added under File Sharing, into its VM. A working directory under `/tmp` fails with a bind-mount error like "not shared from the host", which does not obviously point at the File Sharing setting. Clone or copy these files under your home directory instead.

### 2. Start the stack

```bash
docker compose -f docker-compose.quickstart.yml up -d
```

This starts `ignis-db`, then `ignis-app` once the database reports healthy, then `ignis-reverse-proxy` once the app reports healthy. Only the proxy publishes a host port (`HOST_HTTPS_PORT`, default `443`).

### 3. Seed and verify

Follow [Seeding the database](#seeding-the-database), then [Verifying](#verifying).

### Pinning a version

Both images default to the `latest` published release. To pin a specific one, set `IGNIS_IMAGE_TAG` before starting:

```bash
export IGNIS_IMAGE_TAG=0.2.4-alpha
docker compose -f docker-compose.quickstart.yml --profile seed pull
docker compose -f docker-compose.quickstart.yml up -d
```

!!! info "Tag format"
    Published image tags carry no `v` prefix, even though the Git tags do. Release `v0.2.4-alpha` publishes as `0.2.4-alpha`. Export the variable rather than prefixing a single command, so the app and seed images come from the same release.

---

## Docker from source

Builds `ignis-app` and `ignis-build-db` from this checkout, so it is the only path that picks up local code changes.

### 1. Trust the local CA

This path needs the `caddy` CLI on the host, with `caddy trust` run once. That installs a local CA into your OS and browser trust store, which the proxy then reuses.

```bash
caddy trust
```

### 2. Configure

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

!!! warning "CADDY_DATA_DIR"
    A wrong value silently produces an untrusted certificate. A missing one fails validation instead, reported by name: `required variable CADDY_DATA_DIR is missing a value`.

### 3. Start, seed and verify

```bash
docker compose up -d
```

Then [Seeding the database](#seeding-the-database) and [Verifying](#verifying).

---

## Seeding the database

!!! warning "Required before first use, and destructive"
    A fresh `ignis-db` volume is empty. Seeding drops and recreates all country tables, so it is gated behind the `seed` profile and never runs automatically. Until it has run once, every endpoint that reads the schema will fail.

```bash
<compose prefix> --profile seed run --rm ignis-build-db
```

The TABULA workbook is baked into the `ignis-build-db` image, so there is nothing to download.

## Verifying

```bash
<compose prefix> exec ignis-db psql -U postgres -d ignis -c "\dt tabula.*"
curl -k -s -o /dev/null -w '%{http_code}\n' https://localhost/health
```

The first lists the seeded tables. The second returns `200`. Drop `-k` on the Docker-from-source path, where the certificate chains to the CA you trusted.

## Tearing down

```bash
<compose prefix> down -v
```

The `-v` removes the database volume, so the next start needs seeding again. Omit it to keep the seeded data.

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
    `data/tabula-calculator-lite.xlsx` is stored via Git LFS. Without `git lfs install`, the checkout contains a small text pointer instead of the workbook, and `build_db` fails naming the cause directly: "is a Git LFS pointer, not the workbook itself, run `git lfs install && git lfs pull`".

    ```bash
    git lfs install
    git clone https://github.com/thd-spatial-ai/ignis.git
    cd ignis
    ```

    If you have already cloned without it, run `git lfs install && git lfs pull`. Verify with `ls -l data/tabula-calculator-lite.xlsx`, which should be roughly 11 MB rather than a few hundred bytes.

### Configuration

This path reads `environment/env/app.env` for database settings and `environment/env/common.env` for CORS. The committed values work as they are, apart from these:

| Variable | Change to |
|---|---|
| `DB_HOST` | `localhost`, or wherever your PostgreSQL instance runs |
| `DB_PASSWORD` | your own credential, no default |
| `DB_SSL_MODE` | `disable` only for a local database on the same machine, `require` everywhere else |
| `ALLOWED_ORIGINS` | the browser origins calling ignis directly. Leave unset for server-to-server calls, the intended deployment model |

!!! warning "No API key on this path"
    Running the binary directly bypasses the reverse proxy, so there is no `X-Api-Key` gate and no CORS preflight handling. Bind it to localhost only.

### Build and run

```bash
go build -buildvcs=false -o bin/ ./cmd/...
```

This produces `bin/ignis`, `bin/build_db`, and `bin/validate`.

```bash
./bin/build_db
./bin/ignis
```

`build_db` is destructive: it drops and recreates the `tabula` schema. `ignis` starts on `APP_PORT`, default `8080`.

### Validate

Runs the full 17-level TABULA calculation pipeline against every row in the database and checks that each result stays within 2% of the reference value from the workbook.

```bash
./bin/validate
```

This path only: there is no containerised `validate`. See the [validation report](validation.md) for current results.

---

## Deployment

Use `docker-compose.prod.yml`, which pulls published images and needs no source tree on the target machine. Copy across: the compose file, `.env`, the `env/` directory, and the `caddy/` directory.

### 1. Prepare the `env/` files

!!! danger "Do not deploy the committed env/ files"
    Change all of the following before starting the stack:

    - `POSTGRES_PASSWORD` (`env/db.env`) and `DB_PASSWORD` (`env/app.env`) to a real credential
    - `DB_SSL_MODE` (`env/app.env`) to `require`
    - `IGNIS_API_KEY` (`env/proxy.env`) to a rotated value
    - `ALLOWED_ORIGINS` (`env/common.env`) to the real browser origins, or unset for server-to-server only

### 2. Prepare `.env`

`CADDY_DATA_DIR` is required; `APP_PORT` and `HOST_HTTPS_PORT` default to `8080` and `443`. Set `IGNIS_IMAGE_TAG` to pin a release, which is strongly advised for anything you do not want moving underneath you. See [Pinning a version](#pinning-a-version).

### 3. Set the site address

!!! warning "Set IGNIS_SITE_ADDRESS before deploying"
    It defaults to `localhost`. Set it in `env/proxy.env` to the deployment's real domain, or Caddy will neither serve it nor provision a certificate for it. TLS-ALPN-01 (Caddy's default ACME challenge) works entirely over port 443, which is all the compose file publishes, so no port change is needed.

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

Compose orders startup by health: `ignis-app` waits for a healthy `ignis-db`, and `ignis-reverse-proxy` waits for a healthy `ignis-app`.

`/health` reports process liveness only and does not test the database connection, so a healthy container does not by itself mean the schema is seeded or reachable.