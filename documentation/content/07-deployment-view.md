# Deployment View

## The three containers

ignis runs as three containers, defined in `environment/docker-compose.yml` (development) or `environment/docker-compose.prod.yml` (production — see below). They form the `building-simulation` namespace (the Docker Compose project name), shared with other building-modelling services such as buem. The stack needs nothing on the host except Docker.

1. **ignis-reverse-proxy**: Caddy, from the `caddy:2.11-alpine` image. The only container that publishes a host port. Every request passes through it first.

2. **ignis-app**: the `bin/ignis` binary. In development, built locally from `environment/Dockerfile`; in production, pulled pre-built from `ghcr.io/thd-spatial-ai/ignis` (see "Two ways to run this" below). Listens on the internal port (default 8080), publishes no host port, and reaches the database at `ignis-db`.

3. **ignis-db**: PostgreSQL, from `postgres:17-alpine`. Publishes no host port. Its data lives in a named volume (`ignis-db-data`) that survives `docker compose down`/`up`, but not `down -v`.

Every component is a container, so the whole stack can be built into images, pushed to a registry, and run elsewhere with the same shape of compose file. There is no host-installed database to set up separately.

## Two ways to run this: build from source, or pull a published image

There are two compose files, for two different situations, sharing the same `building-simulation` project name and container names — bringing up either one manages the same logical stack.

**`environment/docker-compose.yml` (development).** `ignis-app` uses `build: context: ./.. dockerfile: environment/Dockerfile` — every `docker compose up` rebuilds `ignis-app` from whatever source is on this machine right now. This needs the repo cloned and is the right choice whenever you're changing ignis's own code. Also defines `ignis-build-db`, a one-off image (see "First-run data load" below).

**`environment/docker-compose.prod.yml` (production).** `ignis-app` instead uses `image: ghcr.io/thd-spatial-ai/ignis:${IGNIS_IMAGE_TAG:-latest}` — nothing is built, the exact image already tested in CI is pulled and run as-is. A machine running this file needs only Docker itself and this repo's small config files (`docker-compose.prod.yml`, `.env`, `docker.env`, `Caddyfile`) — no source tree, no Go toolchain. `ignis-db` and `ignis-reverse-proxy` are unchanged between the two files, since both already pull public images (`postgres:17-alpine`, `caddy:2.11-alpine`) and never needed a build step.

```bash
# production
docker compose -f docker-compose.prod.yml up -d
```

This file is deliberately standalone rather than a partial override merged with `docker-compose.yml`: Compose's merge rules can leave a service with both `build:` and `image:` set, and it can then still fall back to building locally if no matching image is cached — silently defeating the point on a machine with no source code to build from.

**Publishing a new `ignis-app` image**: `.github/workflows/docker-publish.yml` builds and pushes `ghcr.io/thd-spatial-ai/ignis:vX.Y.Z` and `:latest` whenever a `vX.Y.Z` tag is pushed (or the workflow is run manually) — a deliberate release action, not something that happens on every commit. Set `IGNIS_IMAGE_TAG` in `.env` on the deploy machine to pin a specific version rather than always tracking `:latest`. The image is public, so pulling it needs no authentication.

The deployment workflow — release (tag push → CI build → GHCR push) followed by deploy (ordered container startup) — is shown as a figure below. No source code or Go toolchain appears anywhere on the deploy side: the host only ever pulls finished images and runs them. The runtime request path (client → proxy → app → db, gated by `X-Api-Key`) is covered separately below in "No host port on app or database."

## Ports

Two ports, with different rules.

- **Internal port** (`APP_PORT`, default 8080): the port `ignis-app` listens on inside its container. Container isolation means it never clashes with other services, so it stays the same everywhere. It is set once in `.env` and passed to the app, its health check, and the proxy's upstream, rather than hardcoded in each.

- **Host port** (`HOST_HTTPS_PORT`, default 443): the port `ignis-reverse-proxy` publishes. This is the only port that can clash: two services cannot both take host 443. An orchestration layer assigns a free port here per service.

## First-run data load

A fresh `ignis-db` volume is empty. Load the TABULA data once, after first start:

```bash
docker compose --profile seed run --rm ignis-build-db
```

`build_db` is its own image (`environment/Dockerfile.build_db`), separate from `ignis-app` — `ignis-app`'s image contains only the server binary, not `build_db`, `validate`, a compiler, or the source tree. It bakes in a trimmed TABULA workbook (`data/tabula-calculator-lite.xlsx`, just the one sheet this tool reads) rather than reading from a mounted file. It is gated behind Compose's `seed` profile specifically so it never runs as part of a normal `docker compose up`: it prints its own warning on every run ("This will DROP and recreate all country tables!"), so it must only run when explicitly invoked, never automatically on startup. The data persists in the named volume afterward.

## No host port on app or database

`ignis-app` and `ignis-db` declare no `ports:`, so nothing outside the Docker network can connect to them. `ignis-app` is reachable only from `ignis-reverse-proxy`, and `ignis-db` only from `ignis-app`, each by its service name. The proxy is not an add-on in front of an open service: it is the only way in.

## Startup order

`depends_on: condition: service_healthy` chains startup: `ignis-db` must accept connections (checked with `pg_isready`) before `ignis-app` starts, and `ignis-app` must be healthy before `ignis-reverse-proxy` starts. The app's health check makes a real `GET /health` call (not `HEAD`, which the router does not register). Nothing starts serving before what it depends on is ready.

## Certificate trust across recreation

A fresh Caddy container would generate a new, untrusted certificate authority, breaking any trust the browser already had. `ignis-reverse-proxy` avoids this by mounting the host's Caddy data directory (`~/.local/share/caddy`, set via `CADDY_DATA_DIR`) into the container, so it reuses the same CA. A browser that ran `caddy trust` once keeps trusting it.

This is a local-development convenience, tied to one machine. A real deployment replaces it with either public HTTPS (Caddy's built-in Let's Encrypt) or a shared internal CA every client trusts.

## Network

All three containers share one Docker Compose network (`building-simulation_default`). Docker's DNS resolves each service name to its container: the proxy reaches the app at `ignis-app`, the app reaches the database at `ignis-db`. This only works within the same network, which is why keeping everything on it (and nothing on a host port except the proxy) keeps the stack self-contained and closed.
