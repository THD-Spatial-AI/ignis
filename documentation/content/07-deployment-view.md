# Deployment View

## The three containers

ignis runs as three containers, defined in `environment/docker-compose.yml`. They form the `building-simulation` namespace (the Docker Compose project name), shared with other building-modelling services such as buem. The stack needs nothing on the host except Docker.

1. **ignis-reverse-proxy**: Caddy, from the `caddy:2.11-alpine` image. The only container that publishes a host port. Every request passes through it first.

2. **ignis-app**: the `bin/app` binary, built from `environment/Dockerfile`. Listens on the internal port (default 8080), publishes no host port, and reaches the database at `ignis-db`.

3. **ignis-db**: PostgreSQL, from `postgres:17-alpine`. Publishes no host port. Its data lives in a named volume (`ignis-db-data`) that survives `docker compose down`/`up`, but not `down -v`.

Every component is a container, so the whole stack can be built into images, pushed to a registry, and run elsewhere with the same compose file. There is no host-installed database to set up separately.

## Ports

Two ports, with different rules.

- **Internal port** (`APP_PORT`, default 8080): the port `ignis-app` listens on inside its container. Container isolation means it never clashes with other services, so it stays the same everywhere. It is set once in `.env` and passed to the app, its health check, and the proxy's upstream, rather than hardcoded in each.

- **Host port** (`HOST_HTTPS_PORT`, default 443): the port `ignis-reverse-proxy` publishes. This is the only port that can clash: two services cannot both take host 443. An orchestration layer assigns a free port here per service.

## First-run data load

A fresh `ignis-db` volume is empty. Load the TABULA data once, after first start: `docker compose exec ignis-app ./bin/build_db`. The data persists afterward. `build_db` drops and recreates all tables, so it is a manual step, not part of startup.

## No host port on app or database

`ignis-app` and `ignis-db` declare no `ports:`, so nothing outside the Docker network can connect to them. `ignis-app` is reachable only from `ignis-reverse-proxy`, and `ignis-db` only from `ignis-app`, each by its service name. The proxy is not an add-on in front of an open service: it is the only way in.

## Startup order

`depends_on: condition: service_healthy` chains startup: `ignis-db` must accept connections (checked with `pg_isready`) before `ignis-app` starts, and `ignis-app` must be healthy before `ignis-reverse-proxy` starts. The app's health check makes a real `GET /health` call (not `HEAD`, which the router does not register). Nothing starts serving before what it depends on is ready.

## Certificate trust across recreation

A fresh Caddy container would generate a new, untrusted certificate authority, breaking any trust the browser already had. `ignis-reverse-proxy` avoids this by mounting the host's Caddy data directory (`~/.local/share/caddy`, set via `CADDY_DATA_DIR`) into the container, so it reuses the same CA. A browser that ran `caddy trust` once keeps trusting it.

This is a local-development convenience, tied to one machine. A real deployment replaces it with either public HTTPS (Caddy's built-in Let's Encrypt) or a shared internal CA every client trusts.

## Network

All three containers share one Docker Compose network (`building-simulation_default`). Docker's DNS resolves each service name to its container: the proxy reaches the app at `ignis-app`, the app reaches the database at `ignis-db`. This only works within the same network, which is why keeping everything on it (and nothing on a host port except the proxy) keeps the stack self-contained and closed.
