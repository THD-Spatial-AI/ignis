# Getting started

## Clone the repository

```bash
git clone https://github.com/thd-spatial-ai/ignis.git
cd ignis
```

Both setup paths below start from here.

## Docker (recommended)

Runs ignis as three containers — a reverse proxy, the app, and its own PostgreSQL — with nothing else required on the host except Docker itself. The TABULA workbook is baked into the seed image, so there's nothing to download manually.

```bash
cd environment
cp .env.example .env       # APP_PORT, HOST_HTTPS_PORT, CADDY_DATA_DIR
docker compose up -d
```

This starts `ignis-db`, then `ignis-app` once the database is healthy, then `ignis-reverse-proxy` once the app is healthy. Only the proxy publishes a host port (`HOST_HTTPS_PORT`, default `443`).

!!! warning "Load the database before first use"
    A fresh `ignis-db` volume is empty. Seed it once, using the dedicated `seed` profile — this is destructive (drops and recreates all tables), so it never runs automatically:

    ```bash
    docker compose --profile seed run --rm ignis-build-db
    ```

### Quick-start

Quick-start is intended for trying out APIs without any additional prerequisites apart from docker. For a machine that only has Docker (no Go toolchain, `caddy` CLI, `.env` file required). Pulls the pre-built `ignis-app` image from GHCR instead of building from source:

```bash
cd environment
docker compose -f docker-compose.quickstart.yml up -d
cp .env.example .env
docker compose --profile seed run --rm ignis-build-db
```
!!! info "Pinning a version"
    Defaults to the `latest` published image. To pin a specific release, set `IGNIS_IMAGE_TAG` before starting:

    ```bash
    IGNIS_IMAGE_TAG=v0.2.1-alpha docker compose -f docker-compose.quickstart.yml up -d
    
!!! warning "Certificate warning is expected"
    `docker-compose.yml` reuses a `caddy trust`-installed CA from the host so browsers trust `https://localhost` out of the box. This quickstart file skips that step entirely — Caddy's local CA lives in a Docker-managed volume instead, so nothing on the host trusts it. Click through the browser's warning, or pass `-k` / `--no-check-certificate` / "disable SSL verification" from curl, wget, or Postman:

    ```bash
    curl -k -H "X-Api-Key: supersecret123" https://localhost/health
    ```

    Calling the API from Swagger UI's "Try it out" (this site's [API reference](api.md)) needs one extra step: browser JS can't click through a cert warning the way a manual page load can, so a `fetch()` to an untrusted origin just fails outright. Open `https://localhost` directly in a new tab first and click through the warning — most browsers then trust that origin for the rest of the session, and Swagger UI's calls will go through.


    ```

## Manual setup

For local development without containers.

### Prerequisites

| Dependency | Version |
|---|---|
| Go | 1.26+ |
| PostgreSQL | 15 – 17 |
| TABULA Excel workbook | `tabula-calculator.xlsx` |

The Excel workbook is available from [episcope.eu](https://episcope.eu/iee-project/tabula/). Place it in the `data/` directory before running `build_db`.

### Configuration

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | required, no default |
| `DB_NAME` | Database name | `ignis` |
| `DB_SSL_MODE` | TLS mode (`require` / `disable`) | `require` |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins | unset rejects all cross-origin requests |

!!! info "ALLOWED_ORIGINS"
    Only needed when a browser-based client calls ignis directly. For server-to-server calls (the intended deployment model), leave it unset.

!!! warning "DB_SSL_MODE"
    Use `disable` only for local development against a database on the same machine. Set `require` in all other environments.

### Build

```bash
go build -o bin/ignis     cmd/ignis/main.go
go build -o bin/build_db  cmd/build_db/main.go
go build -o bin/validate  cmd/validate/main.go
```

### Load the database

Loads the TABULA workbook into PostgreSQL. This is destructive: it drops and recreates the `tabula` schema.

```bash
./bin/build_db
```

### Run the API

```bash
./bin/ignis   # starts on :8080
```

## Validate

Runs the full 17-level pipeline against every row in the database and checks that the result stays within ±2% of the TABULA reference value. Requires the local Go build above — there is no containerized `validate`.

```bash
./bin/validate
```

See the [validation report](validation.md) for current results.

## Deployment

!!! danger "Do not expose ignis directly to the internet"
    The API has no authentication of its own. Run it behind a reverse proxy on a private network, with no public port exposed on ignis itself.
