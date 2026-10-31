---
audience: developer
---

# Getting started

ignis is an HTTP API that computes TABULA heating demand for building archetypes. It ships as one Docker Compose file, `environment/docker-compose.yml`: the Go application and its own PostgreSQL instance, optionally behind a Caddy reverse proxy that terminates TLS.

## Running with Docker

Run every command from `environment/`. Two choices combine freely: where the images come from, and whether Caddy terminates TLS.

| Images | Command |
|---|---|
| Built from this checkout, picking up local code changes | `docker compose up -d --build` |
| The published release from GHCR, no Go toolchain needed | `docker compose up -d --pull always --no-build` |

Add `--profile https` to either to start Caddy in front of the API (see [https](#https)), or set `COMPOSE_PROFILES=https` in `.env` to make it the default. Without it the stack serves plain HTTP, which is the whole stack where a platform ingress, an orchestrator's proxy or a firewall already terminates TLS.

A third path, [manual setup](#manual-setup), runs the binaries without containers.

!!! warning "Upgrading from `environment/http` and `environment/https`"
    The stack is now the single Compose project `ignis`, replacing `ignis-http` and `ignis-https`. Bring an old stack down from a checkout that still has those directories, or by project name, before the first `up`:

    ```bash
    docker compose -p ignis-http down
    docker compose -p ignis-https down
    ```

    The database volume is now `ignis_db-data`, and `up` seeds it on first start. Remove the old volumes with `docker volume rm ignis-http_db-data ignis-https_db-data` once the new one is populated.

    If `tentacron-net` already exists from `make tentacron-stack`, `up` fails with a label mismatch, because a network made by `docker network create` carries no Compose labels:

    ```
    network tentacron-net was found but has incorrect label
    com.docker.compose.network set to "" (expected: "tentacron-net")
    ```

    An empty value means the network was created outside Compose; remove it with `docker network rm tentacron-net` while nothing is attached, and the first stack up recreates it. Any other value means a compose file declared the network under a different key: Compose compares the key rather than the `name:`, so spell the key `tentacron-net`, the same as the name.

!!! note "Compose warns about the shared network"
    Whichever stack did not create `tentacron-net` prints that it "exists but was not created for project ..." and suggests `external: true`. The warning is expected: one network is shared by several repositories, so every project but the creator reports it. Do not switch it to `external: true`, which requires the network to exist before `up` and so breaks a first run on a clean machine.

## Configuration files

Compose reads an optional `environment/.env`, interpolated on the host, plus the files under `environment/env/`, which are passed into the containers. They are split so each service receives only the variables it reads.

| File | Read by | Holds |
|---|---|---|
| `.env` | Docker Compose, on the host | `COMPOSE_PROFILES`, `IGNIS_IMAGE_TAG`, `HOST_BIND`, `HOST_PORT`, `APP_PORT`, `HOST_HTTPS_PORT`, `CADDY_DATA_DIR` |
| `env/common.env` | api | `ALLOWED_ORIGINS` |
| `env/db.env` | db | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `env/app.env` | api, build-db | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE` |
| `env/proxy.env` | proxy | `IGNIS_SITE_ADDRESS` |

`.env.example` lists every `.env` variable with its default. The `env/` files are committed with local development defaults, including a placeholder database password. Replace them before any deployment: see [Deployment](#deployment).

| Compose service | Image | Role |
|---|---|---|
| `api` | `ghcr.io/thd-spatial-ai/ignis` | HTTP API server |
| `build-db` | `ghcr.io/thd-spatial-ai/ignis-build-db` | TABULA seeder, runs on `up` and exits |
| `db` | `postgres:17-alpine` | PostgreSQL database |
| `proxy` | `caddy:2.11-alpine` | TLS termination, `https` profile only |

Containers are named `ignis-<service>-<n>`, for example `ignis-api-1`. On the shared `tentacron-net` network, callers in other repositories address the API as `ignis`. The bare service name `api` is used by other repositories on that network too, so it is never an address. The alias `ignis-http` still resolves, for callers not yet moved to `ignis`.

## Access control

ignis carries no credential and asks for none. There is no API key on any endpoint.

What limits who can reach it is the port mapping, not `ALLOWED_ORIGINS`. CORS is enforced by a browser on behalf of a page it has loaded; a server-to-server caller sends no `Origin` header and ignores the response headers entirely.

!!! danger "Decide reachability at the network layer"
    The API is published on `HOST_BIND:HOST_PORT`, `127.0.0.1:8088` by default, in both modes, so nothing off the host can connect to it directly. Under `https`, Caddy publishes `HOST_HTTPS_PORT` on every interface. Set `HOST_BIND=0.0.0.0` only where something in front of the host decides who may connect. Deploy ignis only on a network whose access you control, behind a platform that has already authenticated the user.

---

## http

Nothing on the host but Docker. No certificate to trust, no key to send.

```bash
cd environment
docker compose up -d --pull always --no-build
```

The base URL is `http://localhost:8088`. `up` seeds an empty database before it starts the API (see [Seeding the database](#seeding-the-database)), then follow [Verifying](#verifying).

| Variable | Description | Default |
|---|---|---|
| `HOST_BIND` | Host interface the API is published on | `127.0.0.1` |
| `HOST_PORT` | Host port the API is published on | `8088` |
| `APP_PORT` | The API's internal listen port | `8080` |

---

## https

Adds Caddy in front of the API to terminate TLS. Caddy does nothing else: ignis answers CORS preflight itself, and no request carries a credential.

```bash
cd environment
docker compose --profile https up -d --pull always --no-build
```

This starts `db`, then `api` once the database reports healthy, then `proxy` once the API reports healthy. Caddy publishes `HOST_HTTPS_PORT`, default `443`. The base URL is `https://localhost`.

| Variable | Description | Default |
|---|---|---|
| `HOST_HTTPS_PORT` | Host port Caddy publishes | `443` |
| `CADDY_DATA_DIR` | Host directory holding a CA from `caddy trust` | unset: a Docker volume |

### Certificates

How `https://localhost` behaves depends on where Caddy's local CA is stored.

| `CADDY_DATA_DIR` | CA location | Result |
|---|---|---|
| unset | Docker volume `ignis_caddy-data` | Never enters the host trust store. Expect an untrusted-certificate warning. |
| set | The host directory `caddy trust` created | Trusted, no warning, no `-k` needed. |

For a trusted certificate, install the `caddy` CLI on the host, run `caddy trust` once, and set `CADDY_DATA_DIR` in `.env` to the directory it created (Linux `~/.local/share/caddy`, macOS `~/Library/Application Support/Caddy`). A wrong path gives an untrusted certificate, the same as leaving it unset.

Without it, click through the warning in the browser, or pass `-k` (curl), `--no-check-certificate` (wget), or "disable SSL verification" (Postman).

!!! info "Using Swagger UI"
    Calling the API from Swagger UI's "Try it out" on the [API reference](api.md) with an untrusted certificate needs one extra step. Browser JavaScript cannot click through a certificate warning the way a manual page load can, so a `fetch()` to an untrusted origin fails outright. Open `https://localhost` directly in a new tab first and click through the warning. Most browsers then trust that origin for the rest of the session.

!!! warning "Docker Desktop: work from a directory under the home folder"
    Docker Desktop only shares paths under the user's home directory, or another folder added under File Sharing, into its VM. A checkout under `/tmp` fails with a bind-mount error such as "not shared from the host", which does not obviously point at the File Sharing setting.

---

## Seeding the database

The TABULA workbook is baked into the `build-db` image, so there is nothing to download. A seed runs in one transaction: if it fails, nothing is committed and any previous tables stay.

`up` runs `build-db -if-empty` and starts the API only once it exits successfully. An empty database is seeded, a populated one is left unchanged, and a failed seed fails `up` with `service "build-db" didn't complete successfully`.

!!! warning "Rebuilding is destructive"
    Rebuilding a populated database from the workbook drops and recreates all country tables:

    ```bash
    docker compose run --rm build-db -if-empty=false
    ```

!!! warning "Building from source needs Git LFS"
    The workbook is stored in Git LFS and copied into the `build-db` image. A checkout made without `git lfs install` holds a small pointer file instead, and the seed fails. See [Manual setup](#manual-setup).

## Verifying

```bash
docker compose exec db psql -U postgres -d ignis -c "\dt tabula.*"
```

That lists the seeded tables. Then check that the API answers:

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8088/ignis/health   # either mode
curl -s -o /dev/null -w '%{http_code}\n' https://localhost/ignis/health       # https profile
```

Both return `200`. Add `-k` to the second when `CADDY_DATA_DIR` is unset.

## Tearing down

```bash
docker compose --profile https down -v
```

`--profile https` also stops Caddy if it runs, and is harmless otherwise. The `-v` removes the database volume, so the next `up` seeds again. Omit it to keep the seeded data.

## Pinning a version

`--no-build` runs the `latest` published release by default. To pin a specific one, set `IGNIS_IMAGE_TAG` in `.env`, or export it before starting:

```bash
export IGNIS_IMAGE_TAG=0.8.0
docker compose up -d --pull always --no-build
```

!!! info "Tag format"
    Published image tags carry no `v` prefix, even though the Git tags do. Release `v0.8.0` publishes as `0.8.0`. Export the variable rather than prefixing a single command, so the API and seed images come from the same release.

!!! warning "Source builds take the same tag"
    `--build` tags what it builds as `ghcr.io/thd-spatial-ai/ignis:${IGNIS_IMAGE_TAG}`, `latest` by default. Running a release afterwards needs `--pull always`, as above, or the local build keeps running under the release's name.

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

    After a clone made without it, run `git lfs install && git lfs pull`. Verify with `ls -l data/tabula-calculator-lite.xlsx`, which should be roughly 11 MB rather than a few hundred bytes.

### Configuration

This path reads `environment/env/app.env` for database settings and `environment/env/common.env` for CORS. The committed values work as they are, apart from these:

| Variable | Change to |
|---|---|
| `DB_HOST` | `localhost`, or wherever the PostgreSQL instance runs |
| `DB_PASSWORD` | a real credential, no default |
| `DB_SSL_MODE` | `disable` only for a local database on the same machine, `require` everywhere else |
| `ALLOWED_ORIGINS` | the browser origins calling ignis directly. Leave unset for server-to-server calls, the intended deployment model |

!!! warning "Bind it to localhost"
    The binary listens on every interface. Nothing in ignis limits who may call it, so keep it off any interface whose access is not controlled.

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

Runs the full 17-level TABULA calculation pipeline against every row in the database and checks that each result stays within 2.5% of the reference value from the workbook.

```bash
./bin/validate
```

Add `-strict` to exit 1 when any building is outside the tolerance. This path only: there is no containerised `validate`. See the [validation report](validation.md) for current results.

---

## Deployment

Copy the whole `environment/` directory to the target machine: `docker-compose.yml`, `.env`, `env/` and, for `https`, `caddy/`. Release images need no source tree there.

### 1. Prepare the `env/` files

!!! danger "Do not deploy the committed env/ files"
    Change all of the following before starting the stack:

    - `POSTGRES_PASSWORD` (`env/db.env`) and `DB_PASSWORD` (`env/app.env`) to a real credential
    - `DB_SSL_MODE` (`env/app.env`) to `require`
    - `ALLOWED_ORIGINS` (`env/common.env`) to the real browser origins, or unset for server-to-server only

### 2. Prepare `.env`

Set `IGNIS_IMAGE_TAG` to pin a release. This is strongly advised for any deployment that should not move on the next publish. See [Pinning a version](#pinning-a-version).

`HOST_BIND` and `HOST_PORT` default to `127.0.0.1` and `8088`. Set `HOST_BIND=0.0.0.0` only where something in front of the host decides who may connect. For TLS terminated by ignis itself, set `COMPOSE_PROFILES=https`.

### 3. Set the site address (https only)

!!! warning "Set IGNIS_SITE_ADDRESS before deploying"
    It defaults to `localhost`. Set it in `env/proxy.env` to the deployment's real domain, or Caddy will neither serve it nor provision a certificate for it. TLS-ALPN-01 (Caddy's default ACME challenge) works entirely over port 443, which is all Caddy publishes, so no port change is needed.

### 4. Pull and start

`up` seeds an empty database itself.

```bash
docker compose up -d --pull always --no-build
```

!!! warning "Always pull"
    Without `--pull always`, a host that already has the tag cached keeps running that image, with nothing on the host revealing it.

### 5. Verify

```bash
curl -s -o /dev/null -w '%{http_code}\n' https://your-domain/ignis/health
curl -s https://your-domain/api/v1/variants/DE | head -c 200
```

Expect `200`, then a list of variant codes. If the second call fails while health returns `200`, the database is reachable but not seeded.

---

## Health checks

Compose orders startup by health: `api` waits for a healthy `db` and a finished `build-db`, and under `https` the proxy waits for a healthy `api`.

`/ignis/health` reports process liveness only and does not test the database connection, so a healthy container does not by itself mean the schema is seeded or reachable.
