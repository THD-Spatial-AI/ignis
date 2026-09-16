# Getting started

ignis is an HTTP API that computes TABULA heating demand for building archetypes. It ships as a Docker Compose stack: the Go application and its own PostgreSQL instance, optionally behind a Caddy reverse proxy that terminates TLS.

## Choosing an environment

`environment/` holds two self-contained directories. Pick one, `cd` into it, and everything you need is there.

| Directory | Containers | Use when |
|---|---|---|
| [`environment/http`](#http) | app, database | Developing against the API, or deploying where TLS is already terminated upstream |
| [`environment/https`](#https) | app, database, Caddy | ignis has to terminate TLS itself |

TLS is a property of a deployment, not of ignis. Where a platform ingress, an orchestrator's own reverse proxy, or a firewall already handles it, `environment/http` is the whole stack.

Each directory holds one compose file per source of images:

| File | Images |
|---|---|
| `docker-compose.yml` | built from this checkout, so it picks up local code changes |
| `docker-compose.prod.yml` | pulled from GHCR, no Go toolchain or source tree needed |
| `docker-compose.quickstart.yml` (`https` only) | pulled from GHCR, with Caddy's CA in a Docker volume instead of your trust store |

A third path, [manual setup](#manual-setup), runs the binaries without containers.

!!! warning "One stack at a time"
    `environment/http` and `environment/https` use the same container names (`ignis-app`, `ignis-db`), and container names are unique across the host, so the two cannot run side by side. `docker compose down` in one before `up` in the other.

!!! warning "Upgrading from a checkout made before the project rename"
    The Compose project names are now `ignis-http` and `ignis-https`, previously `building-simulation` for both. An existing stack has to come down before the renamed one starts: with it running, `up` fails on the container name rather than replacing it. Remove the old containers by name, which reaches nothing but ignis:

    ```bash
    docker stop ignis-app ignis-db ignis-reverse-proxy
    docker rm ignis-app ignis-db ignis-reverse-proxy
    ```

    Do not use `docker compose -p building-simulation down`. That targets the project rather than this repository, so on a machine where another service still declares the old project name, it removes that service's containers as well, naming neither. The database volume is pinned to its previous name, so it carries over and needs no reseed.

## Configuration files

Each environment reads an optional `.env`, interpolated on the host, plus the files under its own `env/` directory, which are passed into the containers. They are split so each service receives only the variables it reads.

| File | Read by | Holds |
|---|---|---|
| `.env` | Docker Compose, on the host | `APP_PORT`, `IGNIS_IMAGE_TAG`, plus `HOST_BIND`/`HOST_PORT` (http) or `HOST_HTTPS_PORT`/`CADDY_DATA_DIR` (https) |
| `env/common.env` | app | `ALLOWED_ORIGINS` |
| `env/db.env` | db | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `env/app.env` | app, build_db | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE` |
| `env/proxy.env` (https only) | proxy | `IGNIS_SITE_ADDRESS` |

The `env/` files are committed with local development defaults, including a placeholder database password. Replace them before any deployment: see [Deployment](#deployment).

Both dockerfiles stay at `environment/`, shared by the two environments. The published image is the same either way, and CI builds from that one path.

## Access control

ignis carries no credential and asks for none. There is no API key on any endpoint.

What limits who can reach it is the port mapping, not `ALLOWED_ORIGINS`. CORS is enforced by a browser on behalf of a page it has loaded; a server-to-server caller sends no `Origin` header and ignores the response headers entirely.

!!! danger "Decide reachability at the network layer"
    In `environment/http` the app is bound to `127.0.0.1` by default, so nothing off the host can connect. In `environment/https` the app publishes no port at all and only the proxy is exposed. Set `HOST_BIND=0.0.0.0` only where something in front of the host decides who may connect: the intended deployment is an internal network reachable over VPN, behind a platform that has already authenticated the user.

---

## http

Nothing on the host but Docker. No certificate to trust, no key to send.

### 1. Start the stack

```bash
cd environment/http
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

To build from this checkout instead, so local code changes are picked up, drop the `-f` and use `docker compose up -d`.

`.env` is optional: every variable has a default.

| Variable | Description | Default |
|---|---|---|
| `HOST_BIND` | Host interface the app is published on | `127.0.0.1` |
| `HOST_PORT` | Host port the app is published on | `8088` |
| `APP_PORT` | The app's internal listen port | `8080` |

### 2. Seed and verify

Follow [Seeding the database](#seeding-the-database), then [Verifying](#verifying). The base URL is `http://localhost:8088`.

---

## https

Adds Caddy in front of the app to terminate TLS. Caddy does nothing else: ignis answers CORS preflight itself, and no request carries a credential.

### Certificates

How `https://localhost` behaves depends on where Caddy's local CA is stored.

| Compose file | CA location | Result |
|---|---|---|
| `docker-compose.quickstart.yml` | Docker-managed volume | Never enters your trust store. Expect an untrusted-certificate warning. |
| `docker-compose.yml`, `docker-compose.prod.yml` | Host directory (`CADDY_DATA_DIR`), created by `caddy trust` | Trusted, no warning, no `-k` needed. |

On the quickstart path, click through the warning in the browser, or pass `-k` (curl), `--no-check-certificate` (wget), or "disable SSL verification" (Postman).

!!! info "Using Swagger UI"
    Calling the API from Swagger UI's "Try it out" on the [API reference](api.md) needs one extra step on the quickstart path. Browser JavaScript cannot click through a certificate warning the way a manual page load can, so a `fetch()` to an untrusted origin fails outright. Open `https://localhost` directly in a new tab first and click through the warning. Most browsers then trust that origin for the rest of the session.

### Quickstart: pulled images, no host setup

Needs only Docker. No Go toolchain, no `caddy` install, no `.env`.

```bash
git clone https://github.com/thd-spatial-ai/ignis.git
cd ignis/environment/https
docker compose -f docker-compose.quickstart.yml up -d
```

!!! warning "Docker Desktop: work from a directory under your home folder"
    Docker Desktop only shares paths under your home directory, or another folder added under File Sharing, into its VM. A working directory under `/tmp` fails with a bind-mount error like "not shared from the host", which does not obviously point at the File Sharing setting. Clone or copy these files under your home directory instead.

This starts `ignis-db`, then `ignis-app` once the database reports healthy, then `ignis-reverse-proxy` once the app reports healthy. Only the proxy publishes a host port (`HOST_HTTPS_PORT`, default `443`).

Then [Seeding the database](#seeding-the-database) and [Verifying](#verifying).

### From source, with a trusted certificate

Builds `ignis-app` and `ignis-build-db` from this checkout, so it picks up local code changes, and reuses a CA your browser already trusts.

Install the `caddy` CLI on the host and run `caddy trust` once. That installs a local CA into your OS and browser trust store, which the proxy then reuses.

```bash
caddy trust
cd environment/https
cp .env.example .env
docker compose up -d
```

`.env` must define:

| Variable | Description | Default |
|---|---|---|
| `CADDY_DATA_DIR` | Host directory holding the `caddy trust` CA | required, no default |
| `APP_PORT` | The app's internal listen port | `8080` |
| `HOST_HTTPS_PORT` | The proxy's published port | `443` |

!!! warning "CADDY_DATA_DIR"
    A wrong value silently produces an untrusted certificate. A missing one fails validation instead, reported by name: `required variable CADDY_DATA_DIR is missing a value`.

Then [Seeding the database](#seeding-the-database) and [Verifying](#verifying).

---

## Seeding the database

!!! warning "Required before first use, and destructive"
    A fresh `ignis-db` volume is empty. Seeding drops and recreates all country tables, so it is gated behind the `seed` profile and never runs automatically. Until it has run once, every endpoint that reads the schema will fail.

```bash
<compose prefix> --profile seed run --rm ignis-build-db
```

`<compose prefix>` is `docker compose` for `docker-compose.yml`, or `docker compose -f <file>` for the others. On a path that pulls images, `--profile seed` is also needed on the `pull`, since profile-gated services are otherwise skipped.

The TABULA workbook is baked into the `ignis-build-db` image, so there is nothing to download.

## Verifying

```bash
<compose prefix> exec ignis-db psql -U postgres -d ignis -c "\dt tabula.*"
```

That lists the seeded tables. Then check the API answers, using the base URL for your environment:

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8088/ignis/health   # environment/http
curl -s -o /dev/null -w '%{http_code}\n' https://localhost/ignis/health       # environment/https
```

Both return `200`. Add `-k` on the https quickstart path, where the certificate does not chain to a CA you trust.

## Tearing down

```bash
<compose prefix> down -v
```

The `-v` removes the database volume, so the next start needs seeding again. Omit it to keep the seeded data.

## Pinning a version

The paths that pull images default to the `latest` published release. To pin a specific one, set `IGNIS_IMAGE_TAG` before starting:

```bash
export IGNIS_IMAGE_TAG=0.2.4-alpha
<compose prefix> --profile seed pull
<compose prefix> up -d
```

!!! info "Tag format"
    Published image tags carry no `v` prefix, even though the Git tags do. Release `v0.2.4-alpha` publishes as `0.2.4-alpha`. Export the variable rather than prefixing a single command, so the app and seed images come from the same release.

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

This path reads `environment/http/env/app.env` for database settings and `environment/http/env/common.env` for CORS. The committed values work as they are, apart from these:

| Variable | Change to |
|---|---|
| `DB_HOST` | `localhost`, or wherever your PostgreSQL instance runs |
| `DB_PASSWORD` | your own credential, no default |
| `DB_SSL_MODE` | `disable` only for a local database on the same machine, `require` everywhere else |
| `ALLOWED_ORIGINS` | the browser origins calling ignis directly. Leave unset for server-to-server calls, the intended deployment model |

!!! warning "Bind it to localhost"
    The binary listens on every interface. Nothing in ignis limits who may call it, so keep it off any interface you do not control.

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

Use `docker-compose.prod.yml` from whichever environment matches how TLS is handled. It pulls published images and needs no source tree on the target machine.

Copy across the whole directory: the compose file, `.env`, the `env/` directory, and, for `environment/https`, the `caddy/` directory.

!!! info "Service and image names"
    | Compose service | Published image | Role |
    |---|---|---|
    | `ignis-app` | `ghcr.io/thd-spatial-ai/ignis` | HTTP API server |
    | `ignis-build-db` | `ghcr.io/thd-spatial-ai/ignis-build-db` | one-off TABULA seeder, `seed` profile only |
    | `ignis-db` | `postgres:17-alpine` (not built here) | PostgreSQL database |

    The `ignis-app` service publishes without the `-app` suffix. `IGNIS_IMAGE_TAG` pins both `ghcr.io` images to one release.

### 1. Prepare the `env/` files

!!! danger "Do not deploy the committed env/ files"
    Change all of the following before starting the stack:

    - `POSTGRES_PASSWORD` (`env/db.env`) and `DB_PASSWORD` (`env/app.env`) to a real credential
    - `DB_SSL_MODE` (`env/app.env`) to `require`
    - `ALLOWED_ORIGINS` (`env/common.env`) to the real browser origins, or unset for server-to-server only

### 2. Prepare `.env`

For `environment/http`: `HOST_BIND` and `HOST_PORT` default to `127.0.0.1` and `8088`. Set `HOST_BIND=0.0.0.0` only where something in front of the host decides who may connect.

For `environment/https`: `CADDY_DATA_DIR` is required; `APP_PORT` and `HOST_HTTPS_PORT` default to `8080` and `443`.

Set `IGNIS_IMAGE_TAG` to pin a release, which is strongly advised for anything you do not want moving underneath you. See [Pinning a version](#pinning-a-version).

### 3. Set the site address (https only)

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
curl -s -o /dev/null -w '%{http_code}\n' https://your-domain/ignis/health
curl -s https://your-domain/api/v1/variants/DE | head -c 200
```

Expect `200`, then a list of variant codes. If the second call fails while health returns `200`, the database is reachable but not seeded.

---

## Health checks

Compose orders startup by health: `ignis-app` waits for a healthy `ignis-db`, and in `environment/https` the proxy waits for a healthy `ignis-app`.

`/ignis/health` reports process liveness only and does not test the database connection, so a healthy container does not by itself mean the schema is seeded or reachable.
