# Getting started

## Prerequisites

| Dependency | Version |
|---|---|
| Go | 1.26+ |
| PostgreSQL | 15 – 17 |
| TABULA Excel workbook | `tabula-calculator.xlsx` |

The Excel workbook is available from [episcope.eu](https://episcope.eu/iee-project/tabula/). Place it in the `data/` directory before running `build_db`.

## Configuration

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

## Build

```bash
go build -o bin/app       cmd/app/main.go
go build -o bin/build_db  cmd/build_db/main.go
go build -o bin/validate  cmd/validate/main.go
```

## Load the database

Loads the TABULA workbook into PostgreSQL. This is destructive: it drops and recreates the `tabula` schema.

```bash
./bin/build_db
```

## Run the API

```bash
./bin/app   # starts on :8080
```

## Validate

Runs the full 17-level pipeline against every row in the database and checks that the result stays within ±2% of the TABULA reference value.

```bash
./bin/validate
```

See the [validation report](validation.md) for current results.

## Deployment

!!! danger "Do not expose ignis directly to the internet"
    The API has no authentication of its own. Run it behind a reverse proxy on a private network, with no public port exposed on ignis itself.
