# Deployment View

## One compose file, two choices

`environment/docker-compose.yml` is the whole deployment, Compose project `ignis`. Each command chooses where the images come from and whether ignis terminates TLS itself (ADR-007).

| Choice | Flag | Effect |
|---|---|---|
| Images from this checkout | `--build` | Builds `api` and `build-db` from the dockerfiles in `environment/` |
| Images from a release | `--pull always --no-build` | Pulls the tag named by `IGNIS_IMAGE_TAG` from GHCR |
| TLS terminated by ignis | `--profile https` | Adds `proxy` in front of `api` |

Neither path needs anything on the host except Docker. A trusted local certificate additionally needs the `caddy` CLI for its one-off trust step.

The dockerfiles stay at `environment/` because the publishing workflow builds from that path.

## The containers

1. **api**: the `bin/ignis` binary, built from `environment/ignis-app.dockerfile` or pulled as `ghcr.io/thd-spatial-ai/ignis`. Listens on the internal port (default 8080) and reaches the database at `db`.

2. **build-db**: the TABULA seeder, built from `environment/ignis-db.dockerfile` or pulled as `ghcr.io/thd-spatial-ai/ignis-build-db`. Runs on `up`, seeds an empty database and exits; a populated one is left unchanged.

3. **db**: PostgreSQL, from `postgres:17-alpine`. Publishes no host port. Its data lives in the named volume `ignis_db-data`, which survives `docker compose down`/`up`, but not `down -v`.

4. **proxy** (`https` profile only): Caddy, from `caddy:2.11-alpine`. Terminates TLS and forwards plain HTTP to `api`. It holds no CORS configuration and checks no credential.

Containers are named `ignis-<service>-<n>`, for example `ignis-api-1`.

## Ports

- **Internal port** (`APP_PORT`, default 8080): the port `api` listens on inside its container. Container isolation means it never clashes with other services. It is set once and passed to the API, its health check and the proxy's upstream.

- **Host port** (`HOST_PORT`, default 8088, and under `https` also `HOST_HTTPS_PORT`, default 443): the published ports, the only ones that can clash. 8088 avoids 8080, which the EnerPlanET platform's Keycloak binds on every interface.

- **Host interface** (`HOST_BIND`, default `127.0.0.1`): which interface `HOST_PORT` is published on.

## Reachability is a port mapping

`db` declares no `ports:`, so it is reachable only from `api` and `build-db`, by service name.

`api` publishes a port bound to `127.0.0.1` in both modes, so the stack answers only on the host it runs on; under `https`, Caddy is the only way in from another host. Widening `HOST_BIND` to `0.0.0.0` is what exposes ignis to a network, and belongs only where something in front of the host decides who may connect. `ALLOWED_ORIGINS` does not restrict this: CORS is enforced by a browser on behalf of a page, and a server-to-server caller sends no `Origin` header at all.

## Startup order

`depends_on` chains startup: `db` must accept connections (checked with `pg_isready`) before `build-db` runs, `build-db` must exit successfully before `api` starts, and `api` must be healthy before `proxy` starts. The API's health check makes a real `GET /ignis/health` call.

## Certificate trust across recreation

A fresh Caddy container would generate a new, untrusted certificate authority. Setting `CADDY_DATA_DIR` mounts the host directory `caddy trust` created, so the container reuses a CA the browser already trusts. Unset, the CA lives in the Docker volume `ignis_caddy-data`: it survives recreation but is in no trust store, so clients see a certificate warning.

This is a local-development convenience, tied to one machine. A real deployment uses public HTTPS (Caddy's built-in Let's Encrypt), a shared internal CA every client trusts, or plain HTTP behind an ingress that already holds the certificate.

## Networks

The containers share the project's default network, `ignis_default`, where the proxy reaches the API at `api` and the API reaches the database at `db`.

`api` also joins `tentacron-net`, the network other repositories' services share, with the alias `ignis`. That alias is the cross-repository name: it survives any rename of project, service or container. The service name `api` is registered there too, by every repository that uses it, so it resolves to any of them and is never an address. `db` and `build-db` never join `tentacron-net`.
