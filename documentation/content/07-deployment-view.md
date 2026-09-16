# Deployment View

## Two environments, one image

`environment/` holds two self-contained Compose environments. They run the same published image and differ only in whether ignis terminates TLS itself.

| Directory | Containers | Entry point |
|---|---|---|
| `environment/http` | `ignis`, `db` | `ignis`, published on `HOST_BIND:HOST_PORT` |
| `environment/https` | `ignis`, `db`, `proxy` | the proxy, published on `HOST_HTTPS_PORT` |

Each declares its own Docker Compose project, `ignis-http` and `ignis-https`. Both use the same container names, which are unique across the host, so only one runs at a time. Neither needs anything on the host except Docker, and `environment/https` from source additionally needs the `caddy` CLI for its one-off trust step.

The two dockerfiles stay at `environment/` rather than being copied into each directory: the image is identical for both, and the publishing workflow builds from that one path.

## The containers

1. **ignis-app**: the `bin/ignis` binary, built from `environment/ignis-app.dockerfile`. Listens on the internal port (default 8080) and reaches the database at `db`.

2. **db**: PostgreSQL, from `postgres:17-alpine`. Publishes no host port in either environment. Its data lives in a named volume (`ignis-db-data`) that survives `docker compose down`/`up`, but not `down -v`.

3. **proxy** (`environment/https` only): Caddy, from the `caddy:2.11-alpine` image. Terminates TLS and forwards plain HTTP to `ignis`. It holds no CORS configuration and checks no credential.

Every component is a container, so the whole stack can be built into images, pushed to a registry, and run elsewhere with the same compose file. There is no host-installed database to set up separately.

## Ports

- **Internal port** (`APP_PORT`, default 8080): the port `ignis` listens on inside its container. Container isolation means it never clashes with other services, so it stays the same everywhere. It is set once in `.env` and passed to the app, its health check, and, in the HTTPS environment, the proxy's upstream, rather than hardcoded in each.

- **Host port** (`HOST_PORT`, default 8088, or `HOST_HTTPS_PORT`, default 443): the published port. This is the only one that can clash, since two services cannot own the same host port. 8088 avoids 8080, which the EnerPlanET platform's Keycloak binds on every interface. An orchestration layer assigns a free port here per service.

- **Host interface** (`HOST_BIND`, `environment/http` only, default `127.0.0.1`): which interface that port is published on.

## First-run data load

A fresh `db` volume is empty. Load the TABULA data once, after first start, with the `seed` profile: `docker compose --profile seed run --rm build-db`. The data persists afterward. `build_db` drops and recreates all tables, so it is a manual step, not part of startup.

## Reachability is a port mapping

`db` declares no `ports:` in either environment, so it is reachable only from `ignis`, by service name. `ignis` declares none in `environment/https`, where the proxy is the only way in.

In `environment/http` the app publishes a port bound to `127.0.0.1`, so the stack answers only on the host it runs on. Widening that to `0.0.0.0` is what exposes ignis to a network, and belongs only where something in front of the host decides who may connect. `ALLOWED_ORIGINS` does not restrict this: CORS is enforced by a browser on behalf of a page, and a server-to-server caller sends no `Origin` header at all.

## Startup order

`depends_on: condition: service_healthy` chains startup: `db` must accept connections (checked with `pg_isready`) before `ignis` starts, and `ignis` must be healthy before `proxy` starts. The app's health check makes a real `GET /ignis/health` call (not `HEAD`, which the router does not register). Nothing starts serving before what it depends on is ready.

## Certificate trust across recreation

A fresh Caddy container would generate a new, untrusted certificate authority, breaking any trust the browser already had. In `environment/https`, `docker-compose.yml` and `docker-compose.prod.yml` avoid this by mounting the host's Caddy data directory (`~/.local/share/caddy`, set via `CADDY_DATA_DIR`) into the container, so it reuses the same CA. A browser that ran `caddy trust` once keeps trusting it. `docker-compose.quickstart.yml` uses a Docker-managed volume instead, trading the trust step for a certificate warning.

This is a local-development convenience, tied to one machine. A real deployment replaces it with either public HTTPS (Caddy's built-in Let's Encrypt), a shared internal CA every client trusts, or `environment/http` behind an ingress that already holds the certificate.

## Network

All containers in an environment share one Docker Compose network (`ignis-https_default` in the HTTPS environment). The name is derived from the Compose project, so it changes whenever the project does, and services in other repositories have attached to it to resolve `ignis` by name. Renaming the project breaks those callers with a DNS failure at their first outbound call, not at start-up. See ADR-005. Docker's DNS resolves each service name to its container: the proxy reaches the app at `ignis`, the app reaches the database at `db`. This only works within the same network, which is why keeping the database off any host port keeps the stack self-contained.
